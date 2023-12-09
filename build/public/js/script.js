// ---------Responsive-navbar-active-animation-----------
function test(){
	const buttons = document.querySelectorAll('.l_btn');
	$("#navbarSupportedContent").off("click").on("click","li",function(e){
		if (e.target.textContent!=="Выход") {
			$('#navbarSupportedContent ul li').removeClass("active");
			$(this).addClass('active');
			e.preventDefault();
			toggleNavbarLeft(e)
		}
		if (e.target.textContent==="Главная"){
			$( "#wrapper" ).load( "/main_page", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
			});
		} else if (e.target.textContent === "Управление отчётами") {
			buttons.forEach(btn => btn.classList.remove('active'));
			e.target.classList.add('active');
			$( "#wrapper" ).load( "/make_report", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
					}
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.EndTime').addEventListener('change', checkTime)
				document.querySelector('.StartTime').addEventListener('change', checkTime)
				document.querySelector('.proj').addEventListener("change", function(event) {
					const ev_type = "bucket"
					updateDataPage(event, ev_type)
				});
			})
		} else if (e.target.textContent === "Управление тестами") {
			$( "#wrapper" ).load( "/current_tests", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				// document.querySelector('.formWithValidation').addEventListener('submit', handleCreateBucket)
			});
		} else if (e.target.textContent === "Настройки проектов") {
			$("#wrapper").load("/settings", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				updateDataPage(event, "versions")
				updateDataPage(event, "hosts")
				document.getElementById("btn_set_project").addEventListener('click', setActiveProject) // rabotaet
				document.getElementById("btn_set_methodic").addEventListener('click', handleMetodicSet) // rabotaet
				document.getElementById("btn_set_version").addEventListener('click', handleVersionAdd) // rabotaet
				document.getElementById("btn_create_new_bucket").addEventListener('click', handleCreateBucket) // rabotaet
				document.getElementById("btn_set_confl_page").addEventListener('click', handleConflPageAdd)

				// document.getElementById("settings_version_list").addEventListener('change', updateListVersions)

			});
		}
	setTimeout(function() {anim(); }, 200);
	});
}

function anim(e) {
	console.log("ANIM!!")
	let tabsNewAnim = $('#navbarSupportedContent');
	let activeItemNewAnim = $('li.nav-item.active');
	let activeWidthNewAnimHeight = activeItemNewAnim.innerHeight();
	let activeWidthNewAnimWidth = activeItemNewAnim.innerWidth();
	let itemPosNewAnimTop = activeItemNewAnim.position();
	let itemPosNewAnimLeft = activeItemNewAnim.position();
	$(".hori-selector").css({
		"top":itemPosNewAnimTop.top + "px",
		"left":itemPosNewAnimLeft.left + "px",
		"height": activeWidthNewAnimHeight + "px",
		"width": activeWidthNewAnimWidth + "px"
	});
}

function NavbarLeftHandler() {
	const buttons = document.querySelectorAll('.l_btn');
	$('#navbarLeft').off("click").on("click",function(e) {
		let button_text = e.target.textContent
		console.log(button_text);
		buttons.forEach(btn => btn.classList.remove('active'));
		e.target.classList.add('active');
		///////////// Управление отчётами /////////////
		if (button_text === "История отчётов") {
			$( "#wrapper" ).load( "/report_history", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}})
		} else if (button_text === "Создать отчёт") {
			$( "#wrapper" ).load( "/make_report", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.proj').addEventListener("change", function(event) {
					const ev_type = "bucket"
					updateDataPage(event, ev_type)
				});
			})
		} else if (button_text === "Создать бакет") {
			$( "#wrapper" ).load( "/create_bucket", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleCreateBucket)
			});
		} else if (button_text === "Привязка корневой страницы конфлюенс") {
			$( "#wrapper" ).load( "/set_report_homepage", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleSetReportHomepage)
			});
		} else if (button_text === "Привязать методику") {
			$( "#wrapper" ).load( "/set_methodic", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleMetodicSet)
				document.querySelector('.bucke').addEventListener("change", function(event) {
					const ev_type = "projects"
					updateDataPage(event, ev_type)
				});
			})
		} else if (button_text === "Сравнение тестов") {
			$( "#wrapper" ).load( "/compare_release", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleCompareRelease)
				document.querySelector('.bucke').addEventListener("change", function(event) {
					const ev_type = "projects"
					updateDataPage(event, ev_type)
				});
			})
		}
		///////////// Управление тестами /////////////
		else if (button_text === "История тестов") {
			$( "#wrapper" ).load( "/test_history", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}})
		} else if (button_text === "Текущие тесты") {
			$( "#wrapper" ).load( "/current_tests", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				}})
		} else if (button_text === "Запустить тест") {
			$( "#wrapper" ).load( "/start_test", function(responseText, textStatus) {
				if (textStatus === "error") {
					// В случае ошибки загрузки, выводим сообщение
					$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
						"color: red; " +
						"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
				} else {
					document.querySelector('.formWithValidation').addEventListener('submit', handleStartTest)
				}
			})
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
		console.log("TEST");
		NavbarLeftHandler();
	} else {
		navbarLeft.classList.remove('show');
	}
}

function updateDataPage(event, ev_type) {
	let xhr = new XMLHttpRequest();
	let data = {};
	let url
	var select
	if (ev_type === "bucket") {
		console.log(`update bucket!`)
		data["project"] = document.querySelector('#project_options').value
		data["StartTime"] = document.querySelector('.StartTime').value
		data["EndTime"] = document.querySelector('.EndTime').value
		console.log(`project: `, data["project"])
		url = "/get_project_buckets"
		select = $('#bucket_options')
	} else if (ev_type === "versions") {
		console.log(`update versions!`)
		data["project"] = document.querySelector('#settings_activeproject').value
		console.log(`project: `, data["project"])
		url = "/beeload/get/version"
		select = $('#settings_version_list')
	} else if (ev_type === "projects") {
		console.log(`update projects!`)
		data["project"] = document.querySelector('#bucket_options').value
		console.log(`project: `, data["project"])
		url = "/get_bucket_projects"
		select = $('#project_options')
	} else if (ev_type === "hosts") {
		console.log(`update hosts!`)
		url = "/get_host_list"
		select = $('#settings_host')
	}
	xhr.open("POST", url, true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.onreadystatechange = function() {
		if (this.readyState === XMLHttpRequest.DONE && this.status === 200) {
			let response = JSON.parse(this.responseText);
			console.log(response);
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

function handleStartTest(event) {
	let form = document.querySelector('.formWithValidation')
	let git_url = form.querySelector('.gitUrl')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["gitlab"] = git_url.value
	data["count"] = 1
	data["resource"] = "2cpu4ram"
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "http://172.17.155.33:9999/create", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
}

function handleMetodicSet(event) {
	let form = document.querySelector('.settings_page');
	let bucket = form.querySelector('.bucket');
	let page = form.querySelector('.page');
	let version = form.querySelector('.version');
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault();
	let pageNumber = page.value.replace(/\D/g, ''); // Оставляет только цифры
	page.value = pageNumber;

	data["bucket"] = bucket.value;
	data["version"] = version.value;
	data["page"] = page.value;
	console.log("JSON: ", JSON.stringify(data));

	//TODO: Реализация

	xhr.open("POST", "/beeload/add/methodic", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = document.querySelector('.error_message')
	msg.textContent = 'Статус: Методика привязана';
}

function handleVersionAdd(event) {
	let form = document.querySelector('.settings_page')
	let version = form.querySelector('#set_version')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["version"] = version.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/add/version", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = document.querySelector('.error_message')
	msg.textContent = 'Статус: Новая версия создана'
}

function handleConflPageAdd(event) {
	let form = document.querySelector('.settings_page')
	let page = form.querySelector('#confl_page_url')
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["page"] = page.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/add/confl_page", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = document.querySelector('.error_message')
	msg.textContent = 'Статус: Новая страница привязана'

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
	let msg = document.querySelector('.error_message')
	let xhr = new XMLHttpRequest();
	let data = {};
	// Выполняем действие с выбранным значением, например, отправляем на сервер
	// alert("Выбран проект: " + selectedValue);
	// msg.textContent = "Выбран проект: "+ selectedValue
	updateDataPage(event, "versions")
	data["project_name"] = selectedValue
	xhr.open("POST", "/beeload/set/project", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.addEventListener("load", function() {
		if (xhr.status === 200) { // Коды ответов
			msg.textContent = "Выбран проект: "+ selectedValue
			// alert("Успешный запрос!");
		} else if (xhr.status === 404) {
			msg.textContent = "Ресурс не найден: " + xhr.status
			alert("Ресурс не найден.");
		} else if (xhr.status === 500) {
			msg.textContent = "Внутренняя ошибка сервера: " + xhr.status
			// alert("Внутренняя ошибка сервера.");
		} else {
			msg.textContent = "Неизвестный код ответа: " + xhr.status
			// alert("Неизвестный код ответа: " + xhr.status);
		}
	});
	xhr.send(JSON.stringify(data));
}

function handleCreateBucket(event) {
	event.preventDefault();
	let host =document.querySelector('#settings_host')
	let bucket = document.querySelector('#new_bucket_name');
	console.log("handleCreateBucket")
	let xhr = new XMLHttpRequest();
	let data = {};
	event.preventDefault()
	data["host"] = host.value
	data["bucket"] = bucket.value
	console.log("JSON: ", JSON.stringify(data))
	xhr.open("POST", "/beeload/create/bucket", true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.send(JSON.stringify(data));
	let msg = document.querySelector('.error_message')
	msg.textContent = "Статус: Новый бакет " + data["bucket"] + " создан"
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

function checkTime() {
	let form = document.querySelector('.formWithValidation')
	let startTimeInput = form.querySelector('.StartTime')
	let endTimeInput = form.querySelector('.EndTime')
	const startTime = new Date(startTimeInput.value);
	const endTime = new Date(endTimeInput.value);

	if (startTime > endTime) {
		alert("Время начала не может быть больше времени завершения.");
		startTimeInput.value = ""; // Очищаем поле времени начала
	}
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

		} else if  (this.status !== 200) {

		}
	}
	xhr.send(JSON.stringify(data));
}

function toTimestamp(strDate){
	let datum = Date.parse(strDate);
	return datum/1000;
}

$(document).ready(function(){
	$("#wrapper").load("/main_page", function(responseText, textStatus, jqXHR) {
		if (textStatus === "error") {
			// В случае ошибки загрузки, выводим сообщение
			$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
				"color: red; " +
				"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
		}
	});
	setTimeout(function() {anim(); }, 200);
	setTimeout(function(){ test(); }, 100);
	$(window).on('resize', function() {
		setTimeout(function(){ anim(); }, 500);
	});
	$(".navbar-toggler").click(function(){
		$(".navbar-collapse").slideToggle(300);
		setTimeout(function(){ test(); });
	});
});

// --------------add active class-on another-page move----------
// jQuery(document).ready(function($){
// 	let path = window.location.pathname.split("/").pop();
// 	if (path == '') {
// 		path = '/';
// 	}
// 	console.log("path = " + path);
// 	let target = $('#mainNavLink');
// 	// Add active class to target link
// 	target.parent().addClass('active');
// });

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
