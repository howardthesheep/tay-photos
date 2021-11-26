import ApiManager from "./apiManager.js";

let manager = new ApiManager();

window.onload = () => {
	loadPicOnclick();
}
window.shareImage = shareImage;
window.downloadImage = downloadImage;

function loadPicOnclick() {
	document.querySelectorAll(".photo-gallery img").forEach((image) => {
		image.addEventListener('click', (evt) => {
			console.log("photo clicked");
			let imgId = evt.target.getAttribute("data-id");
			location.href = manager.getPhotoLink(imgId);
		});
	});
}

export function shareImage() {
	alert('Share image');
}

export function downloadImage() {
	alert('Download image');
}