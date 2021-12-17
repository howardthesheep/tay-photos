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
	//let galleryPhotoJson = await galleryPhotos.json();

	updateGalleryInfo(galleryInfoJson);
	//updateGalleryPhotos(galleryPhotoJson);
	// TODO: Populate page with Gallery Photos
}

function updateGalleryInfo(json) {
	let date = new Date(json.create_time);

	document.querySelector("#gallery-name").textContent = json.name;
	document.querySelector("#gallery-date").textContent = date.toDateString();
}