import ApiManager from "./apiManager";

const api = new ApiManager();

document.onload = async () => {
	// Based on the gallery we are accessing, localhost:6969/gallery.html?id=1
	// Retrieve the photos to display on page, providing a gallery ID

	// Populate page with header photo, collections, & collection photos
	await loadGalleryData();
}

// Loads all of the dynamic gallery data and populates the page with it
async function loadGalleryData() {
	let galleryInfo, galleryPhotos;

	try {
		galleryInfo = await api.getGalleryInfo();
		galleryPhotos = await api.getGalleryPhotos();
	} catch (e) {
		throw(e);
	}

	console.log(galleryInfo);
	console.log(galleryPhotos);
	// TODO: Populate page with GalleryInfo & Photos
}