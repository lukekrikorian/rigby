function setLoggedIn(state){
	window.isLoggedIn = state;
	localStorage.setItem("isLoggedIn", state.toString())
}

function isLoggedIn(){
	if (window.isLoggedIn === undefined) {
		const val = localStorage.getItem("isLoggedIn");
		if (val === null) {
			localStorage.setItem("isLoggedIn", "false")
			window.isLoggedIn = false;
			return false;
		} else {
			return val === "true";
		}
	}
	return window.isLoggedIn;
}

function getSelf(){
	fetch("/api/users/me", { credentials: "same-origin" })
	.then(resp => resp.json())
	.then(resp => { window.user = resp })
	.catch(console.error);
}

export { setLoggedIn, isLoggedIn, getSelf };