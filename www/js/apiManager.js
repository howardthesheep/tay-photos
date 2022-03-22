export default class ApiManager {
	apiLocation;

	constructor() {
		this.apiLocation = "http://localhost:6969";
	}

	getPhotoLink(photoId) {
		// TODO: This will need to be removed when done testing
		if (photoId == null) {
			return `${this.apiLocation}/photo.html`;
		}

		return(`${this.apiLocation}/photo/${photoId}`);
	}

	// Gets all the gallery information that is not photos or collection names
	async getGalleryInfo(galleryId) {
		try {
			return await this._apiRequest(`${this.apiLocation}/gallery?id=${galleryId}`, 'GET')
		} catch (e) {
			throw(e);
		}
	}

	// Gets everything related to gallery photos & collection names
	async getGalleryPhotos(galleryId) {
		try {
			return await this._apiRequest(`${this.apiLocation}/gallery/photos?id=${galleryId}`, 'GET')
		} catch (e) {
			throw(e);
		}
	}

	// Sends a user & password combination to backend
	// Returns: an apiToken to be stored in sessionStorage
	async login(userData) {
		const dataStr = JSON.stringify(userData);
		try {
			return await this._apiRequest(`${this.apiLocation}/user/login`, 'POST', dataStr);
		} catch (e) {
			throw(e);
		}
	}

	async deleteUser(userData) {
		const dataStr = JSON.stringify(userData);
		try {
			return await this._apiRequest(`${this.apiLocation}/user`, 'DELETE', dataStr)
		} catch (e) {
			throw(e);
		}
	}

	async createUser(userData) {
		const dataStr = JSON.stringify(userData);
		try {
			return await this._apiRequest(`${this.apiLocation}/user`, 'POST', dataStr);
		} catch (e) {
			throw(e);
		}
	}

	async updateUser(userData) {
		const dataStr = JSON.stringify(userData);
		try {
			return await this._apiRequest(`${this.apiLocation}/user`, 'PUT', dataStr);
		} catch (e) {
			throw(e);
		}
	}

	// Removes the apiToken in storage and forces user to re-auth
	async logout() {
		try {
			return await this._apiRequest(`${this.apiLocation}/user/logout`, 'POST', '')
		} catch (e) {
			throw(e);
		}
	}

	// Helper function which does all the heavy lifting of creating, configuring, and sending requests to backend
	async _apiRequest(endpoint, requestMethod, body) {
		return new Promise((fufill, reject) => {
			const requestData = {
				method: requestMethod,
				mode: 'cors',
				cache: 'default',
				body: body == null ? null : body,
			};
			const request = new Request(endpoint, requestData);

			fetch(request).then((response) => {
				fufill(response);
			}).catch((err) => {
				console.error(`Error during API Request: ${err}\n\nRequest: ${request}`)
				reject();
			})
		});
	}
}