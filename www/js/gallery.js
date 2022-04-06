import ApiManager from "./apiManager.js";

const api = new ApiManager();

// Populate page with header photo, collections, & collection photos
await loadGalleryData();

// Disable right clicking on the page, preventing ~most~ user's from saving
// images easily. Although any dev worth their salt can figure out how to
// download them
disableRightClick();

// Loads all of the dynamic gallery data and populates the page with it
async function loadGalleryData() {
	let galleryInfo, galleryPhotos;

	// Parse url args
	let rawQueryArgs = location.search;
	let queryParams = new URLSearchParams(rawQueryArgs);

	let id = queryParams.get("id");

	try {
		galleryInfo = await api.getGalleryInfo(id);
		galleryPhotos = await api.getGalleryPhotos(id);
	} catch (e) {
		throw(e);
	}

	let galleryInfoJson = await galleryInfo.json();
	let galleryPhotoJson = await galleryPhotos.json();

	updateGalleryInfo(galleryInfoJson);
	updateGalleryPhotos(galleryPhotoJson);
}

// Dynamically updates page information with data from backend
function updateGalleryInfo(json) {
	let date = new Date(json.create_time);

	document.querySelector("#gallery-name").textContent = json.name;
	document.querySelector("#gallery-date").textContent = date.toDateString();
}

// Dynamically populates the page with gallery section and associated images
function updateGalleryPhotos(json) {
	const collections = parsePhotoCollections(json);

	for (const collection of collections) {
		// Create section header
		let sectionTitle = document.createElement("h2");
		sectionTitle.textContent = collection.name;

		// Create section itself
		let section = document.createElement("section");
		section.classList.add("photo-gallery");

		// Add img tags to the new section
		for (const photoId of collection.photoIds) {
			const img = document.createElement("img");
			img.src = api.apiLocation + "/gallery/photo?id=" + photoId
			img.alt = "Gallery Image"
			img.onclick = lightboxImage;
			section.append(img)
		}

		// Append everything to body
		document.body.append(sectionTitle);
		document.body.append(section);
	}
}

// Shows an image in a fullscreen view overlayed on top of the current page
function lightboxImage(evt) {
	const imageSource = evt.target.src;

	// Create our image itself
	const lightboxImage = document.createElement("img");
	lightboxImage.src = imageSource;
	lightboxImage.alt = "Gallery Image";
	lightboxImage.onclick = removeLightbox;
	lightboxImage.classList.add("lightboxImage");

	// Add Share button
	const shareBtn = document.createElement("button");
	shareBtn.textContent = "Share";
	shareBtn.onclick = (evt) => {}; // TODO: Implement Sharing feature
	shareBtn.classList.add("lightboxButton");

	// Add download button
	const downloadBtn = document.createElement("button");
	downloadBtn.textContent = "Download";
	downloadBtn.onclick = (evt) => {}; // TODO: Implement Downloading feature
	downloadBtn.classList.add("lightboxButton");

	// Create Buttons container
	const buttonContainer = document.createElement("div");
	buttonContainer.classList.add("lightboxButtonContainer");
	buttonContainer.append(shareBtn, downloadBtn);

	// Create image container
	const imageContainer = document.createElement("div");
	imageContainer.classList.add("lightboxContainer");
	imageContainer.onclick = removeLightbox;
	imageContainer.append(lightboxImage);
	imageContainer.append(buttonContainer);

	document.body.append(imageContainer);
}

// Removes all .lightboxContainer from the DOM tree
function removeLightbox(evt) {
	console.log(evt.target.nodeName);
	if (evt.target.nodeName.toLowerCase() === "button") return;

	const containers = document.querySelectorAll(".lightboxContainer");
	for (const container of containers) {
		container.remove();
	}
}

// Disable right-clicking anywhere on the <body> or on its descendants
function disableRightClick() {
	document.body.oncontextmenu = () => {return false};
}

// Parses json into groups of array of json objects of the form
// {name: <collectionName>, photoIds: [<array of photoIds in this collection>]}
function parsePhotoCollections(json) {
	const collectionObj = [];

	let collections = [];
	for (const jsonElement of json) {
		collections.push(jsonElement.collection);
	}
	const uniqueCollections = collections.filter(function (value, index, self) {
		return self.indexOf(value) === index;
	})

	// TODO: mapping could probably be more efficient
	let newEntry;
	for (const uniqueCollection of uniqueCollections) {
		newEntry = {
			name: uniqueCollection,
			photoIds: []
		}

		for (const obj of json) {
			if (obj.collection === uniqueCollection) {
				newEntry.photoIds.push(obj.id);
			}
		}
		collectionObj.push(newEntry)
	}

	return collectionObj;
}