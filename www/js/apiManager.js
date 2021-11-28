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

		return(`${this.apiLocation}/photo/${photoId}`);
	}

	// TODO
	// Sends a user & hashed password combination to backend
	// Returns: a JWT to be stored in sessionStorage
	login(userData) {

	}

	// TODO
	// Removes the JWT in sessionStorage and forces user to re-auth
	logout() {

	}
}