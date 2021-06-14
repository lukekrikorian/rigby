var form = document.getElementsByTagName("form")[0],
		errorText = document.getElementById("errorText");

form.addEventListener("submit", function (e) {
	e.preventDefault();
	
	var data = {};
	form.childNodes.forEach(function (child) {
		if (!child.type || child.type === "submit") return;
		if (child.type === "checkbox") {
			data[child.name] = child.checked ? 1 : 0;
			return;
		}
		data[child.name] = child.value;
	});

	Array.from(form.attributes).forEach(function (attr) {
		if (attr.nodeName === "method" || attr.nodeName === "action") return;
		data[attr.nodeName] = attr.nodeValue;
	});

	console.log(data);
	
	var req = new XMLHttpRequest();
	req.open(form.method, form.action);
	req.setRequestHeader("Content-Type", "application/json");
	req.addEventListener("load", function(){
		if (req.status === 303) {
			location = req.responseText;
		} else if (req.status !== 200) {
			errorText.textContent = req.responseText;
		}
	});
	req.send(JSON.stringify(data));

});
