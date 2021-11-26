export default class ApiManager {
	apiLocation;


	constructor() {
		this.apiLocation = "http://localhost:6969";
	}

	getPhotoLink(photoId) {
		// TODO: This will need to be removed when done testing
		if (photoId == null) {
			return `${this.apiLocation}/html/photo.html`;
		}

		return(`${this.apiLocation}/photo/${photoId}`)
	}
}