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
	shareBtn.onclick = showShareDialog;
	shareBtn.classList.add("lightboxButton");

	// Add download button
	const downloadBtn = document.createElement("button");
	downloadBtn.textContent = "Download";
	downloadBtn.onclick = downloadImage;
	downloadBtn.classList.add("lightboxButton");

	// Create Buttons container
	const buttonContainer = document.createElement("div");
	buttonContainer.classList.add("lightboxButtonContainer");
	buttonContainer.append(shareBtn, downloadBtn);

	// Create container that fits size of image
	const imgContainer = document.createElement("div");
	imgContainer.classList.add("lightboxImageContainer");
	imgContainer.append(lightboxImage);
	imgContainer.append(buttonContainer);

	// Create container of everything
	const container = document.createElement("div");
	container.classList.add("lightboxContainer");
	container.onclick = removeLightbox;
	container.append(imgContainer);

	document.body.append(container);
}

function showShareDialog(evt) {

	// Create Social Buttons
	const twitterBtn = document.createElement("img");
	twitterBtn.classList.add("socialButton");
	const instagramBtn = document.createElement("img");
	instagramBtn.classList.add("socialButton");
	const socialLinkContainer = document.createElement("div");
	socialLinkContainer.classList.add("socialLinkContainer");
	socialLinkContainer.appendChild(twitterBtn);
	socialLinkContainer.appendChild(instagramBtn);

	// Create shareable link & copy button
	const linkText = document.createElement("input");
	linkText.id = "shareLink";
	linkText.inputMode = "text";
	// TODO: this should include an anchor to the image being shared, then have a check which opens
	// 		 the image in lightbox if it's id is present in the url anchor
	linkText.value = "http://localhost:6969/gallery.html?id=1f579937-d99c-4e0f-98a7-6a7e0976bcab";
	linkText.readOnly = true;
	const copyButton = document.createElement("button");
	copyButton.innerText = "Copy";
	copyButton.onclick = copyShareLink;
	const copyLinkContainer = document.createElement("div");
	copyLinkContainer.appendChild(linkText);
	copyLinkContainer.appendChild(copyButton);

	// Create dialog header
	const dialogX = document.createElement("button"); // TODO: replace w/ icon
	dialogX.id = "closeShareBtn"
	dialogX.onclick = closeShareDialog;
	dialogX.innerText = "X";
	const dialogH3 = document.createElement("text");
	dialogH3.innerText = "Share";
	const dialogHeader = document.createElement("h3");
	dialogHeader.classList.add("shareDialogHeader");
	dialogHeader.appendChild(dialogH3);
	dialogHeader.appendChild(dialogX);

	// Create dialog itself
	const shareDialog = document.createElement("div");
	shareDialog.classList.add("shareDialog");
	shareDialog.appendChild(dialogHeader);
	shareDialog.appendChild(socialLinkContainer);
	shareDialog.appendChild(copyLinkContainer);

	// Create root container
	const shareContainer = document.createElement("section");
	shareContainer.classList.add("shareDialogContainer");
	shareContainer.onclick = closeShareDialog;
	shareContainer.appendChild(shareDialog);

	document.body.appendChild(shareContainer);
}

function closeShareDialog(evt) {
	const closeAreaClicked = evt.target.classList.contains("shareDialogContainer");
	const closeBtnClicked = evt.target.id === "closeShareBtn";
	if (closeAreaClicked || closeBtnClicked) {
		const container = document.querySelector(".shareDialogContainer");
		container?.remove();
	}
}

// Copies the share link to the clipboard from the open share dialog
function copyShareLink(_) {
	const shareLink = document.querySelector("#shareLink").value;

	if (!navigator.clipboard) {
		try {
			document.execCommand('copy', null, shareLink);
		} catch (err) {
			console.error('Fallback: Oops, unable to copy', err);
		}
		return;
	}
	navigator.clipboard.writeText(shareLink).then(null, function(err) {
		console.error('Async: Could not copy text: ', err);
	});
}

// Downloads an image through the browser
function downloadImage(_) {
	const link = document.createElement("a");
	const imgElement = document.querySelector(".lightboxImage");
	link.href = imgElement.src;
	link.download = imgElement.src.split("=").pop() + ".jpeg";
	document.body.appendChild(link);
	link.click();
	document.body.removeChild(link);
}

// Removes all .lightboxContainer from the DOM tree
function removeLightbox(evt) {
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