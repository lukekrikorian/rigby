var voteButton = document.getElementById("voteButton");

if (!/session/i.test(document.cookie)) {
	voteButton.parentElement.classList.add("dark");
	voteButton.parentElement.innerHTML = voteButton.textContent;
}

voteButton.addEventListener("click", function(e){
	e.preventDefault();
	var req = new XMLHttpRequest();
	req.open("POST", voteButton.href);
	req.addEventListener("load", function(){
		if (req.status === 200) {
			voteButton.textContent = req.responseText + " Votes";
		}
	});
	req.send();
});