import ApiManager from "./apiManager.js";

const api = new ApiManager();

// Populate page with header photo, collections, & collection photos
await loadGalleryData();

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

function updateGalleryInfo(json) {
	let date = new Date(json.create_time);

	document.querySelector("#gallery-name").textContent = json.name;
	document.querySelector("#gallery-date").textContent = date.toDateString();
}

function updateGalleryPhotos(json) {
	const collections = parsePhotoCollections(json);

	console.log(collections);

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
			section.append(img)
		}

		// Append everything to body
		document.body.append(sectionTitle);
		document.body.append(section);
	}
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