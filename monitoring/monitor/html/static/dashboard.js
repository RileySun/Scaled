//Functions
function show(e) {
	e.removeEventListener("animationend", animationEnd)
	e.style.display = "block"
	e.style.animation = "show 0.5s linear 1 forwards";
}

function hide(e) {
	e.style.animation = "hide 0.5s linear 1 forwards";
	e.addEventListener("animationend", animationEnd, false)
}

function animationEnd(e) {
	e.target.style.display = "none";
}



async function restart(e) {
	const id = e.target.parentElement.parentElement.getAttribute("data-id")
	const url = window.location + "restart/"+id;
	
	try {
		const response = await fetch(url);
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		
		const raw = await response.text();
		switch(raw) {
			case "OK":
				window.location.reload();
			case "Error":
				//TODO Notify
			default:
				
		}
	}
	catch (error) {
		console.error(error.message);
	}
}

async function shutdown(e) {
	const id = e.target.parentElement.parentElement.getAttribute("data-id")
	const url = window.location + "shutdown/"+id;
	
	try {
		const response = await fetch(url);
		if (!response.ok) {
			throw new Error(`Response status: ${response.status}`);
		}
		
		const raw = await response.text();
		switch(raw) {
			case "OK":
				window.location.reload();
			case "Error":
				//TODO Notify
			default:
				
		}
	}
	catch (error) {
		console.error(error.message);
	}
}

//Load
window.addEventListener("load", () => {
	const elems = document.getElementsByClassName("Item")
	const arr = Array.from(elems)
	
	arr.forEach((e) => {
		const health = e.getElementsByClassName("Item-Health")[0];
		const hidden = e.getElementsByClassName("Item-Hidden")[0];
		const rest = e.getElementsByClassName("Item-Control")[0];
		const shut = e.getElementsByClassName("Item-Control")[1];
		
		health.addEventListener("mouseover", () => {show(hidden)}, false)
		health.addEventListener("mouseout", () => {hide(hidden)}, false)
		rest.addEventListener("click", restart, false)
		shut.addEventListener("click", shutdown, false)
	}); 
	
}, false);