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

	testRequest() {
		const initData = {
			method: 'POST',
			mode: 'cors',
			cache: 'default'
		};
		const request = new Request(`${this.apiLocation}/photo/test`, initData);

		fetch(request)
			.then((response) => {
				console.log(response);
			}).catch((err) => {
				console.error(`Error during testRequest: ${err}`);
		});
	}

	testUserCRUD() {
		let initData = {
			method: 'post',
		}

		// Create User
		let request = new Request(`${this.apiLocation}/user/`, initData);
		fetch(request)
			.then((response) => {
				console.log(response);
			}).catch((err) => {
			console.error(`Error during testRequest: ${err}`);
		});

		// Update User
		initData.method = 'PUT';
		request = new Request(`${this.apiLocation}/user`, initData);
		fetch(request)
			.then((response) => {
				console.log(response);
			}).catch((err) => {
			console.error(`Error during testRequest: ${err}`);
		});

		// Get User
		initData.method = 'GET';
		request = new Request(`${this.apiLocation}/user`, initData);
		fetch(request)
			.then((response) => {
				console.log(response);
			}).catch((err) => {
			console.error(`Error during testRequest: ${err}`);
		});

		// Delete User
		initData.method = 'DELETE';
		request = new Request(`${this.apiLocation}/user`, initData);
		fetch(request)
			.then((response) => {
				console.log(response);
			}).catch((err) => {
			console.error(`Error during testRequest: ${err}`);
		});
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