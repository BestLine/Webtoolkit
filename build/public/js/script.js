// ---------Responsive-navbar-active-animation-----------
function test(){
	let tabsNewAnim = $('#navbarSupportedContent');
	let selectorNewAnim = $('#navbarSupportedContent').find('li').length;
	let activeItemNewAnim = tabsNewAnim.find('.active');
	let activeWidthNewAnimHeight = activeItemNewAnim.innerHeight();
	let activeWidthNewAnimWidth = activeItemNewAnim.innerWidth();
	let itemPosNewAnimTop = activeItemNewAnim.position();
	let itemPosNewAnimLeft = activeItemNewAnim.position();
	const buttons = document.querySelectorAll('.l_btn');
	$(".hori-selector").css({
		"top":itemPosNewAnimTop.top + "px",
		"left":itemPosNewAnimLeft.left + "px",
		"height": activeWidthNewAnimHeight + "px",
		"width": activeWidthNewAnimWidth + "px"
	});
	$("#navbarSupportedContent").on("click","li",function(e){
		$('#navbarSupportedContent ul li').removeClass("active");
		if (e.target.textContent!=="Выход") {
			$(this).addClass('active');
			e.preventDefault();
			toggleNavbarLeft(e)
		}
		let activeWidthNewAnimHeight = $(this).innerHeight();
		let activeWidthNewAnimWidth = $(this).innerWidth();
		let itemPosNewAnimTop = $(this).position();
		let itemPosNewAnimLeft = $(this).position();

		// toggleNavbarLeft(e)
		if (e.target.textContent==="Главная"){
			$( "#wrapper" ).load( "/main_page" );
		} else if (e.target.textContent === "Управление отчётами") {
			buttons.forEach(btn => btn.classList.remove('active'));
			e.target.classList.add('active');
			$( "#wrapper" ).load( "/make_report", function() {
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.proj').addEventListener("change", updateListBuckets)
			})
		} else if (e.target.textContent === "Управление тестами") {
			$( "#wrapper" ).load( "/current_tests", function () {
				// document.querySelector('.formWithValidation').addEventListener('submit', handleCreateBucket)
			});
		} else if (e.target.textContent === "Настройки проектов") {
			$("#wrapper").load("/settings", function () {
				updateListVersions()
				document.getElementById("btn_set_project").addEventListener('click', setActiveProject)
				document.getElementById("btn_set_methodic").addEventListener('click', handleMetodicSet)
				document.getElementById("btn_set_version").addEventListener('click', handleVersionAdd)
				// document.getElementById("settings_version_list").addEventListener('change', updateListVersions)

			});
		}

		$(".hori-selector").css({
			"top":itemPosNewAnimTop.top + "px",
			"left":itemPosNewAnimLeft.left + "px",
			"height": activeWidthNewAnimHeight + "px",
			"width": activeWidthNewAnimWidth + "px"
		});
	});
}

function NavbarLeftHandler() {
	const buttons = document.querySelectorAll('.l_btn');
	$('#navbarLeft').on("click",function(e) {
		let button_text = e.target.textContent
		console.log(button_text);
		buttons.forEach(btn => btn.classList.remove('active'));
		e.target.classList.add('active');
		///////////// Управление отчётами /////////////
		if (button_text === "История отчётов") {
			$( "#wrapper" ).load( "/report_history")
		} else if (button_text === "Создать отчёт") {
			$( "#wrapper" ).load( "/make_report", function() {
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.proj').addEventListener("change", updateListBuckets)
			})
		} else if (button_text === "Создать бакет") {
			$( "#wrapper" ).load( "/create_bucket", function () {
				document.querySelector('.formWithValidation').addEventListener('submit', handleCreateBucket)
			});
		} else if (button_text === "Привязка корневой страницы конфлюенс") {
			$( "#wrapper" ).load( "/set_report_homepage", function () {
				document.querySelector('.formWithValidation').addEventListener('submit', handleSetReportHomepage)
			});
		} else if (button_text === "Привязать методику") {
			$( "#wrapper" ).load( "/set_methodic", function() {
				document.querySelector('.formWithValidation').addEventListener('submit', handleMetodicSet)
				document.querySelector('.bucke').addEventListener("change", updateListProjects)
			})
		} else if (button_text === "Сравнение тестов") {
			$( "#wrapper" ).load( "/compare_release", function() {
				document.querySelector('.formWithValidation').addEventListener('submit', handleCompareRelease)
				document.querySelector('.bucke').addEventListener("change", updateListProjects)
			})
		}
		///////////// Управление тестами /////////////
		else if (button_text === "История тестов") {
			$( "#wrapper" ).load( "/test_history")
		} else if (button_text === "Текущие тесты") {
			$( "#wrapper" ).load( "/current_tests")
		} else if (button_text === "Запустить тест") {
			$( "#wrapper" ).load( "/start_test")
		}
	});
}

function toggleNavbarLeft(e) {
	var navbarLeft = document.getElementById('navbarLeft');
	navbarLeft.innerHTML = ''; // Очищаем содержимое навигационного блока слева

	// Создаем новый блок с кнопками при нажатии на любую кнопку, кроме первой
	if (e.target.textContent !== 'Главная' && e.target.textContent !== 'Настройки проектов') {
		var newButtons = '';
		if (e.target.textContent === 'Управление отчётами') {
			newButtons = '<button class="l_btn">История отчётов</button>' +
				'<button class="l_btn active">Создать отчёт</button>' +
				'<button class="l_btn">Сравнение тестов</button>';
		} else if (e.target.textContent === 'Управление тестами') {
			newButtons = '<button class="l_btn">История тестов</button>' +
				'<button class="l_btn active">Текущие тесты</button>' +
				'<button class="l_btn">Запустить тест</button>';
		}

		navbarLeft.innerHTML = newButtons;
		if (!navbarLeft.classList.contains('show')) {
			navbarLeft.classList.add('show'); // Добавляем класс 'show' для показа блока с анимацией
		}
		NavbarLeftHandler();
	} else {
		navbarLeft.classList.remove('show');
	}
}

function updateListBuckets(event) {
	console.log(`update!!`)
	let xhr = new XMLHttpRequest();
	let data = {};
	let new_options = [];
	data["project"] = document.querySelector('#project_options').value
	data["StartTime"] = document.querySelector('.StartTime').value
	data["EndTime"] = document.querySelector('.EndTime').value

	console.log(`project: `, document.querySelector('#project_options').value)
	xhr.open("POST", "/get_project_buckets", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function() {
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = JSON.parse(this.responseText);
			console.log(response);
			var select = $('#bucket_options')
			select.empty();
			for (var j = 0; j < response.length; j++){
				console.log(response[j]);
				select.append("<option value='" +response[j]+ "'>" +response[j]+ "     </option>");
			}
		}
	}
	xhr.send(JSON.stringify(data));
}

function  updateListVersions(event) {
	console.log(`update versions!`)
	let xhr = new XMLHttpRequest();
	let data = {};
	let new_versions = [];
	data["project"] = document.querySelector('#settings_activeproject').value
	xhr.open("POST", "/get_version_list", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function() {
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = JSON.parse(this.responseText);
			console.log(response);
			var select = $('#settings_version_list')
			select.empty();
			for (var j = 0; j < response.length; j++){
				console.log(response[j]);
				select.append("<option value='" +response[j]+ "'>" +response[j]+ "     </option>");
			}
		}
	}
	xhr.send(JSON.stringify(data));
}

function updateListProjects(event) {
	console.log(`update proj!`)
	let xhr = new XMLHttpRequest();
	let data = {};
	let new_options = [];
	data["project"] = document.querySelector('#bucket_options').value

	console.log(`project: `, document.querySelector('#bucket_options').value)
	xhr.open("POST", "/get_bucket_projects", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function() {
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = JSON.parse(this.responseText);
			console.log(response);
			var select = $('#project_options')
			select.empty();
			for (var j = 0; j < response.length; j++){
				console.log(response[j]);
				select.append("<option value='" +response[j]+ "'>" +response[j]+ "     </option>");
			}
		}
	}
	xhr.send(JSON.stringify(data));
}

function handleCompareRelease(event) {
	let form = document.querySelector('.formWithValidation')
	let project = form.querySelector('.project')
	let bucket = form.querySelector('.bucket')
	let version = form.querySelector('.version')
	let defect = form.querySelector('.defect')
	let test_type = form.querySelector('.test_type')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["project"] = project.value
	data["bucket"] = bucket.value
	data["version"] = version.value
	data["defect"] = defect.checked
	data["test_type"] = test_type.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/compare/release", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
}

function handleMetodicSet(event) {
	let form = document.querySelector('.settings_page')
	let bucket = form.querySelector('.bucket')
	let page = form.querySelector('.page')
	let version = form.querySelector('.version')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["bucket"] = bucket.value
	data["version"] = version.value
	data["page"] = page.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/add/methodic", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = form.querySelector('.msg')
	msg.textContent = 'методика привязана'
}

function handleVersionAdd(event) {
	let form = document.querySelector('.settings_page')
	let version = form.querySelector('#set_version')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["version"] = version.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "beeload/add/version", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = form.querySelector('.msg')
	msg.textContent = 'новая версия создана'
}

function handleSetReportHomepage(event) {
	let form = document.querySelector('.formWithValidation')
	let project = form.querySelector('.page')
	let bucket = form.querySelector('.bucket')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["id"] =  parseInt(project.value)
	data["bucket"] = bucket.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/add/home", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function (message) {
		var responceMsg = document.querySelector('.error_message');
		var anim_border = document.querySelector('.form_wrapper');
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = this.responseText;
			console.log(response);
			responceMsg.textContent = response;
			responceMsg.style.color = '#33cc33';
			responceMsg.style.border = '3px solid #33cc33';
			responceMsg.style.animation="errorAnim .2s forwards";
			responceMsg.style.textAlign = 'center';
			responceMsg.style.fontSize = '24px';
			anim_border.style.border = '0px solid #33cc33';
			anim_border.style.transition = 'border-width 0.5s';
			anim_border.style.borderWidth = '5px';
			anim_border.classList.add('animation');
		} else if  (this.status !== 200) {
			responceMsg.textContent = "[ERROR] Код ответа сервера: "+ this.status +this.responseText;
			responceMsg.style.color = '#F65656';
			responceMsg.style.textAlign = 'center';
			responceMsg.style.fontSize = '24px';
			anim_border.style.border = '0px solid #F65656';
			anim_border.style.transition = 'border-width 0.5s';
			anim_border.style.borderWidth = '5px';
			anim_border.classList.add('animation');
			responceMsg.style.border = '3px solid #F65656';
			responceMsg.style.animation="errorAnim .2s forwards";
		}
	}
	xhr.send(JSON.stringify(data));
}

function setActiveProject(event) {
	var select = document.getElementById("settings_activeproject");
	var selectedValue = select.options[select.selectedIndex].value;

	// Выполняем действие с выбранным значением, например, отправляем на сервер
	alert("Выбран проект: " + selectedValue);
	let data = {
		project: selectedValue
	};
	fetch('/beeload/set/project', {
		method: 'POST',
			headers: {
			'Content-Type': 'application/json; charset=UTF-8'
		},
		body: JSON.stringify(data)
	})
	.then(response => {
		if (response.ok) {
			if (response.redirected) {
				window.location.href = response.url;
			} else {
				let response = this.responseText;
				console.log(response);
				responceMsg.textContent = response;
				responceMsg.style.color = '#33cc33';
				responceMsg.style.border = '3px solid #33cc33';
				responceMsg.style.animation = "errorAnim .2s forwards";
				responceMsg.style.textAlign = 'center';
				responceMsg.style.fontSize = '24px';
				anim_border.style.border = '0px solid #33cc33';
				anim_border.style.transition = 'border-width 0.5s';
				anim_border.style.borderWidth = '5px';
				anim_border.classList.add('animation');
			}
		} else {
			responceMsg.textContent = "[ERROR] Код ответа сервера: " + this.status + this.responseText;
			responceMsg.style.color = '#F65656';
			responceMsg.style.textAlign = 'center';
			responceMsg.style.fontSize = '24px';
			anim_border.style.border = '0px solid #F65656';
			anim_border.style.transition = 'border-width 0.5s';
			anim_border.style.borderWidth = '5px';
			anim_border.classList.add('animation');
			responceMsg.style.border = '3px solid #F65656';
			responceMsg.style.animation = "errorAnim .2s forwards";
		}
	})
		.catch(error => {
			// Обработка ошибок сети или других проблем
		});

}

function handleCreateBucket(event) {
	event.preventDefault();
	let form = document.querySelector('.formWithValidation');
	let host = form.querySelector('.project');
	let bucket = form.querySelector('.bucket');
	let data = {
		host: host.value,
		bucket: bucket.value
	};

	fetch('/beeload/create/bucket', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json; charset=UTF-8'
		},
		body: JSON.stringify(data)
	})
		.then(response => {
			if (response.ok) {
				if (response.redirected) {
					window.location.href = response.url;
				} else {
					let response = this.responseText;
					console.log(response);
					responceMsg.textContent = response;
					responceMsg.style.color = '#33cc33';
					responceMsg.style.border = '3px solid #33cc33';
					responceMsg.style.animation = "errorAnim .2s forwards";
					responceMsg.style.textAlign = 'center';
					responceMsg.style.fontSize = '24px';
					anim_border.style.border = '0px solid #33cc33';
					anim_border.style.transition = 'border-width 0.5s';
					anim_border.style.borderWidth = '5px';
					anim_border.classList.add('animation');
				}
			} else {
				responceMsg.textContent = "[ERROR] Код ответа сервера: " + this.status + this.responseText;
				responceMsg.style.color = '#F65656';
				responceMsg.style.textAlign = 'center';
				responceMsg.style.fontSize = '24px';
				anim_border.style.border = '0px solid #F65656';
				anim_border.style.transition = 'border-width 0.5s';
				anim_border.style.borderWidth = '5px';
				anim_border.classList.add('animation');
				responceMsg.style.border = '3px solid #F65656';
				responceMsg.style.animation = "errorAnim .2s forwards";
			}
		})
		.catch(error => {
			// Обработка ошибок сети или других проблем
		});
}


function handleCompare(event) {
	let form = document.querySelector('.formWithValidation')
	let bucket = form.querySelector('.bucket')
	let StartTime = form.querySelector('.StartTime')
	let EndTime = form.querySelector('.EndTime')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["bucket"] = bucket.value
	data["EndTime"] = EndTime.value
	data["StartTime"] = StartTime.value
	data["start_timestamp"] = toTimestamp(StartTime.value)
	data["end_timestamp"] = toTimestamp(EndTime.value)
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/test_make_compare", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
}

function handleMakeReport(event) {
	let form = document.querySelector('.formWithValidation')
	let bucket = form.querySelector('#bucket_options')
	let project = form.querySelector('.project')
	let StartTime = form.querySelector('.StartTime')
	let EndTime = form.querySelector('.EndTime')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	StartTime.classList.add("invalid")
	data["application"] = bucket.value
	data["bucket"] = project.value
	data["EndTime"] = EndTime.value
	data["StartTime"] = StartTime.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/make/report", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function (message) {
		var responceMsg = document.querySelector('.error_message');
		var anim_border = document.querySelector('.form_wrapper');
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = JSON.parse(this.responseText);
			console.log(response);
			responceMsg.textContent = response;
			responceMsg.style.color = '#06a142';
			responceMsg.style.border = '3px solid ##06a142';
			responceMsg.style.animation="errorAnim .2s forwards";
			responceMsg.style.textAlign = 'center';
			responceMsg.style.fontSize = '24px';
			anim_border.style.border = '0px solid #33cc33';
			anim_border.style.transition = 'border-width 0.5s';
			anim_border.style.borderWidth = '5px';
			anim_border.classList.add('animation');
		} else if  (this.status !== 200) {
			responceMsg.textContent = "[ERROR] Код ответа сервера: "+ this.status;
			responceMsg.style.color = '#F65656';
			responceMsg.style.textAlign = 'center';
			responceMsg.style.fontSize = '24px';
			anim_border.style.border = '0px solid #F65656';
			anim_border.style.transition = 'border-width 0.5s';
			anim_border.style.borderWidth = '5px';
			anim_border.classList.add('animation');
			responceMsg.style.border = '3px solid #F65656';
			responceMsg.style.animation="errorAnim .2s forwards";
		}
	}
	xhr.send(JSON.stringify(data));
}

function toTimestamp(strDate){
	let datum = Date.parse(strDate);
	return datum/1000;
}

$(document).ready(function($){
	setTimeout(function(){ test(); });
});
$(window).on('resize', function(){
	setTimeout(function(){ test(); }, 500);
});
$(".navbar-toggler").click(function(){
	$(".navbar-collapse").slideToggle(300);
	setTimeout(function(){ test(); });
});

// --------------add active class-on another-page move----------
$(window).on('DOMContentLoaded',function () {
	// Get current path and find target link
	let path = window.location.pathname.split("/").pop();

	// Account for home page with empty path
	if (path == '') {
		path = '/';
	}
	console.log("path = " + path);

	let target = $('#navbarSupportedContent ul li a[href="' + path + '"]');
	// Add active class to target link
	target.parent().addClass('active');

	// Проверяем, был ли уже выполнен .load()
	if (!$("#wrapper").data("loaded")) {
		$("#wrapper").data("loaded", true);

		$( "#wrapper" ).load( "/main_page" );
	}
});




// Add active class on another page linked
// ==========================================
// $(window).on('load',function () {
//     var current = location.pathname;
//     console.log("current path = "+current);
//     $('#navbarSupportedContent ul li a').each(function(){
//         var $this = $(this);
//         // if the current path is like this link, make it active
//         if($this.attr('href').indexOf(current) !== -1){
//             $this.parent().addClass('active');
//             $this.parents('.menu-submenu').addClass('show-dropdown');
//             $this.parents('.menu-submenu').parent().addClass('active');
//         }else{
//             $this.parent().removeClass('active');
//         }
//     })
// });



// document.querySelector('#a_background').addEventListener("change", function() {
// 	let xhr = new XMLHttpRequest();
// 	let data = {};
// 	let new_options = [];
// 	data["project"] = this.value
// 	xhr.open("POST", "/get_project_buckets", true);
// 	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
// 	xhr.onreadystatechange = function() {
// 		if (xhr.readyState === XMLHttpRequest.DONE) {
// 			new_options = xhr.responseText;
// 		}
// 	}
// 	xhr.send(JSON.stringify(data));
// 	$('yourSelectList').empty();
// 	$each(new_options, function(value) {
// 		new Element('option')
// 			.set('text', value)
// 			.inject($('yourSelectList'));
// 	});
// });
