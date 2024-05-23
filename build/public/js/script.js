// ---------Responsive-navbar-active-animation-----------
function test(){
	const buttons = document.querySelectorAll('.l_btn');
	let wrapper = $( "#wrapper" )
	$("#navbarSupportedContent").off("click").on("click","li",function(e){
		if (e.target.textContent!=="Выход") {
			$('#navbarSupportedContent ul li').removeClass("active");
			$(this).addClass('active');
			e.preventDefault();
			toggleNavbarLeft(e)
		}
		if (e.target.textContent==="Главная"){
			wrapper.html("")
			wrapper.load( "/main_page", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}
			});
		} else if (e.target.textContent === "Управление отчётами") {
			buttons.forEach(btn => btn.classList.remove('active'));
			e.target.classList.add('active');
			wrapper.html("")
			wrapper.load( "/make_report", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
					}
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.EndTime').addEventListener('change', checkTime)
				document.querySelector('.StartTime').addEventListener('change', checkTime)
				document.querySelector('.proj').addEventListener("change", function(event) {
					// const ev_type = "bucket"
					const ev_type = "test"
					console.log("UPDATE : ", ev_type)
					updateDataPage(event, ev_type)
					//TODO: добавить обработку списка тестов
					//TODO: добавить обработку списка тестов
					//TODO: добавить обработку списка тестов
					//TODO: добавить обработку списка тестов
				});
			})
		} else if (e.target.textContent === "Управление тестами") {
			wrapper.html("")
			wrapper.load( "/tests", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}
				// document.querySelector('.formWithValidation').addEventListener('submit', handleCreateBucket)
			});
		} else if (e.target.textContent === "Настройки проектов") {
			wrapper.html("")
			wrapper.load("/settings", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}
				// updateDataPage(event, "versions")
				updateDataPage(event, "hosts")
				document.getElementById("btn_set_project").addEventListener('click', setActiveProject) // rabotaet
				// document.getElementById("btn_set_methodic").addEventListener('click', handleMetodicSet) // rabotaet
				// document.getElementById("btn_set_version").addEventListener('click', handleVersionAdd) // rabotaet
				// document.getElementById("btn_create_new_bucket").addEventListener('click', handleCreateBucket) // rabotaet
				// document.getElementById("btn_set_confl_page").addEventListener('click', handleConflPageAdd)

				// document.getElementById("settings_version_list").addEventListener('change', updateListVersions)

			});
		}
	setTimeout(function() {anim(); }, 200);
	});
}

function show_error() {
	$("#wrapper").html("<div class=\"main_page\" style=\"text-align: center; " +
		"color: red; " +
		"font-size: 20px;\">Ошибка при загрузке содержимого. Сервер недоступен.</div>");
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
		let wrapper = $( "#wrapper" )
		console.log(button_text);
		buttons.forEach(btn => btn.classList.remove('active'));
		e.target.classList.add('active');
		///////////// Управление отчётами /////////////
		if (button_text === "История отчётов") {
			wrapper.html("")
			wrapper.load( "/report_history", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}})
		} else if (button_text === "Создать отчёт") {
			wrapper.html("")
			wrapper.load( "/make_report", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleMakeReport)
				document.querySelector('.proj').addEventListener("change", function(event) {
					const ev_type = "projects"
					updateDataPage(event, ev_type)
				});
			})
		}  else if (button_text === "Сравнение тестов") {
			wrapper.html("")
			wrapper.load( "/compare_release", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}
				document.querySelector('.formWithValidation').addEventListener('submit', handleCompareRelease)
				document.querySelector('.bucke').addEventListener("change", function(event) {
					const ev_type = "projects"
					updateDataPage(event, ev_type)
				});
			})
		}
		///////////// Управление тестами /////////////
		else if (button_text === "Тесты") {
			wrapper.html("")
			wrapper.load( "/tests", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				}})
		} else if (button_text === "Запустить тест") {
			wrapper.html("")
			wrapper.load( "/start_test", function(responseText, textStatus) {
				if (textStatus === "error") {
					show_error()
				} else {
					// document.querySelector('.formWithValidation').addEventListener('submit', handleStartTestParseEnv)
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
			newButtons =
				'<button class="l_btn active">Тесты</button>' +
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

function updateDataPage(event, ev_type) {
	let xhr = new XMLHttpRequest();
	let data = {};
	let url
	var select
	if (ev_type === "bucket") {
		console.log(`update bucket!`)
		data["project"] = document.querySelector('#project_options').value
		// data["StartTime"] = document.querySelector('.StartTime').value
		// data["EndTime"] = document.querySelector('.EndTime').value
		console.log(`project: `, data["project"])
		url = "/get_project_buckets"
		select = $('#bucket_options')
	} else if (ev_type === "test"){
		console.log(`update test!`)
		data["project"] = document.querySelector('#project_options').value
		data["StartTime"] = document.querySelector('.StartTime').value + ":00Z";
		data["EndTime"] = document.querySelector('.EndTime').value + ":00Z";
		console.log(`project: `, data["project"])
		url = "/get_list_of_tests"
		select = $('#test_options')
	} else if (ev_type === "versions") {
		console.log(`update versions!`)
		data["project"] = document.querySelector('#settings_activeproject').value
		console.log(`project: `, data["project"])
		url = "/beeload/get/version"
		select = $('#settings_version_list')
	} else if (ev_type === "projects") {
		console.log(`update projects!`)
		data["project"] = document.querySelector('#project_options').value
		console.log(`Project options: `, data["project"])
		url = "/get_bucket_projects"
		select = $('#test_options')
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

function goBack() {
	$( "#wrapper" ).load( "/start_test", function(responseText, textStatus) {
		if (textStatus === "error") {
			show_error()
		} else {
			document.querySelector('.formWithValidation').addEventListener('click', handleStartTestParseEnv)
		}
	})
}

function startTest() {
	event.preventDefault()
	let form = document.getElementById('TestStartForm');
	let gitlab = document.getElementById('url').value;
	let count = parseInt(document.getElementById('quantity').value, 10);
	let resource = document.getElementsByName('generator')[0].value;
	let filename = document.getElementsByName('filename')[0].value;

	let envsData = [];
	let envFields = document.querySelectorAll('.Envs div');

	if (!form.checkValidity()) {
		form.reportValidity();
		return;
	}

	envFields.forEach(function(envField) {
		let key = envField.querySelector('.area_label').textContent.toLowerCase();
		let value = envField.querySelector('input').value;
		envsData.push({ key: key, value: value });
	});

	let data = {
		gitlab: gitlab,
		count: count,
		resource: resource,
		data: envsData,
		testplan: filename
	};

	console.log("startTest JSON: ", JSON.stringify(data))
	send_request_with_notification2(data, '/beeload/test/create', "Статус: Запрос на запуск теста отправлен")
}

function handleStartTestParseEnv(event, scenario) {
	console.log("handleStartTestParseEnv scenario: ", scenario)
	let form = document.querySelector('.formWithValidation')
	let git_url = form.querySelector('.gitUrl')
	let data = {};
	let url = ''
	event.preventDefault()
	data["gitlab"] = git_url.value
	console.log("handleStartTestParseEnv JSON: ", JSON.stringify(data))
	if (scenario===true) {
		url = '/parse/env'
	} else {
		url = '/parse/env/custom'
	}
	fetch(url, {
		method: 'POST',
		body: JSON.stringify(data)
	})
		.then(response => response.text())
		.then(result => $("#wrapper").html("<div class=\"main_page\">" + result + "</div>"))
}

function handleStartTest(event) {
	let form = document.querySelector('.formWithValidation')
	let git_url = form.querySelector('.gitUrl')
	let count = form.querySelector('.genCount')
	let gen_type = form.querySelector('.genType')
	let data = {};
	event.preventDefault()

	if (git_url.value.trim() === '' || count.value.trim() === '' || gen_type.value.trim() === '') {
		alert("Необходимо заполнить все поля формы!");
		return;
	}

	data["gitlab"] = git_url.value
	data["count"] = count.value
	data["resource"] = gen_type.value
	console.log("handleStartTest JSON: ", JSON.stringify(data))
	send_request_with_notification(data, "/beeload/test/create", "Статус: Запрос на запуск теста отправлен")
}

function setActiveProject(event) {
	var select = document.getElementById("settings_activeproject");
	var selectedValue = select.options[select.selectedIndex].value;
	let msg = document.querySelector('.error_message')
	let xhr = new XMLHttpRequest();
	let data = {};
	data["project_name"] = selectedValue
	send_request_with_notification(data, "/beeload/set/project", "Выбран проект: " + selectedValue)
}

function send_request_with_notification2(data, url, message) {
	console.log("send_request_with_notification");
	let msg = document.querySelector('.error_message');

	fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json; charset=UTF-8',
		},
		body: JSON.stringify(data),
	})
		.then(response => {
			if (!response.ok) {
				throw new Error('Network response was not ok');
			}
			return response.json();
		})
		.then(result => {
			// Дополнительная логика после успешного запроса
			msg.textContent = message;
		})
		.catch(error => {
			// Обработка ошибок
			if (error instanceof TypeError) {
				msg.textContent = "Сетевая ошибка: " + error.message;
			} else {
				msg.textContent = "Неизвестная ошибка: " + error.message;
			}
		});
}

function send_request_with_notification(data, url, message) {
	// В данный момент может только в пост
	console.log("send_request_with_notification")
	let xhr = new XMLHttpRequest();
	let msg = document.querySelector('.error_message')
	xhr.open("POST", url, true);
	xhr.setRequestHeader('Content-Type', 'application/json; charset=UTF-8');
	xhr.addEventListener("load", function() {
		if (xhr.status === 200) { // Коды ответов
			msg.textContent = message
		} else if (xhr.status === 404) {
			msg.textContent = "Ресурс не найден: " + xhr.status
		} else if (xhr.status === 500) {
			msg.textContent = "Внутренняя ошибка сервера: " + xhr.status
		} else {
			msg.textContent = "Неизвестный код ответа: " + xhr.status
		}
	});
	xhr.send(JSON.stringify(data));
}

function handleCompare(event) {
	let form = document.querySelector('.formWithValidation')
	let bucket = form.querySelector('.bucket')
	let StartTime = form.querySelector('.StartTime')
	let EndTime = form.querySelector('.EndTime')
	let data = {};
	event.preventDefault()
	data["bucket"] = bucket.value
	data["EndTime"] = EndTime.value
	data["StartTime"] = StartTime.value
	data["start_timestamp"] = toTimestamp(StartTime.value)
	data["end_timestamp"] = toTimestamp(EndTime.value)
	console.log("JSON: ", JSON.stringify(data))
	send_request_with_notification(data, "/test_make_compare", "Статус: сравнение тест")
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
	let project = form.querySelector('#project_options')
	let test = form.querySelector('#test_options')
	let StartTime = form.querySelector('.StartTime')
	let EndTime = form.querySelector('.EndTime')
	let data = {};
	event.preventDefault()
	StartTime.classList.add("invalid")
	data["application"] = test.value
	data["bucket"] = project.value
	data["EndTime"] = EndTime.value
	data["StartTime"] = StartTime.value
	console.log("JSON: ", JSON.stringify(data))
	send_request_with_notification(data, "/beeload/make/report", "Статус: Генерация отчёта запрошена.")
}

function toTimestamp(strDate){
	let datum = Date.parse(strDate);
	return datum/1000;
}

function openNewTab(project, test) {
	var url =
		"http://grafana.qa-auto.vimpelcom.ru/d/ms5SULTnz/apache-jmeter-dashboard-using-core-influxdbbackendlistenerclient?" +
		"orgId=1&" +
		"refresh=30s&" +
		"var-data_source=" + project + "&" +
		"var-application=" + test + "&" +
		"var-measurement_name=jmeter&" +
		"var-send_interval=5&" +
		"from=now-15m&to=now";
	window.open(url, '_blank');
}

$(document).ready(function(){
	$("#wrapper").load("/main_page", function(responseText, textStatus, jqXHR) {
		if (textStatus === "error") {
			show_error()
		}
	});
	setTimeout(function() {anim();}, 200);
	setTimeout(function(){test();}, 100);
	$(window).on('resize', function() {
		setTimeout(function(){anim();}, 500);
	});
	$(".navbar-toggler").click(function(){
		$(".navbar-collapse").slideToggle(300);
		setTimeout(function(){test();});
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
