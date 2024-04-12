var sessionID = "0"
var windowFocus = true

window.onresize = function() {
	scanElementsSize();
}

window.onbeforeunload = function(event) {
	sendMessage( "session-close{session=" + sessionID +"}" );
}

window.onblur = function(event) {
	windowFocus = false
	sendMessage( "session-pause{session=" + sessionID +"}" );
}

function sessionInfo() {

	const touch_screen = (('ontouchstart' in document.documentElement) || (navigator.maxTouchPoints > 0) || (navigator.msMaxTouchPoints > 0)) ? "1" : "0";
	var message = "startSession{touch=" + touch_screen 
	
	const style = window.getComputedStyle(document.body);
	if (style) {
		var direction = style.getPropertyValue('direction');
		if (direction) {
			message += ",direction=" + direction
		}
	}

	const lang = window.navigator.language;
	if (lang) {
		message += ",language=\"" + lang + "\"";
	}

	const langs = window.navigator.languages;
	if (langs) {
		message += ",languages=\"" + langs + "\"";
	}

	const userAgent = window.navigator.userAgent
	if (userAgent) {
		message += ",user-agent=\"" + userAgent + "\"";
	}

	const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
	if (darkThemeMq.matches) {
		message += ",dark=1";
	} 

	const pixelRatio = window.devicePixelRatio;
	if (pixelRatio) {
		message += ",pixel-ratio=" + pixelRatio;
	}

	if (localStorage.length > 0) {
		message += ",storage="
		lead = "_{"
		for (var i = 0; i < localStorage.length; i++) {
			var key = localStorage.key(i)
			var value = localStorage.getItem(key)
			key = key.replaceAll(/\\/g, "\\\\")
			key = key.replaceAll(/\"/g, "\\\"")
			key = key.replaceAll(/\'/g, "\\\'")
			value = value.replaceAll(/\\/g, "\\\\")
			value = value.replaceAll(/\"/g, "\\\"")
			value = value.replaceAll(/\'/g, "\\\'")
			message += lead + "\"" + key + "\"=\"" + value + "\""
			lead = ","
		}
		message += "}"
	}

	return message + "}";
}

function restartSession() {
	sendMessage( sessionInfo() );
}

function getIntAttribute(element, tag) {
	let value = element.getAttribute(tag);
	if (value) {
		return value;
	}
	return 0;
}

function scanElementsSize() {
	var element = document.getElementById("ruiRootView");
	if (element) {
		let rect = element.getBoundingClientRect();
		let width = getIntAttribute(element, "data-width");
		let height = getIntAttribute(element, "data-height");
		if (rect.width > 0 && rect.height > 0 && (width != rect.width || height != rect.height)) {
			element.setAttribute("data-width", rect.width);
			element.setAttribute("data-height", rect.height);
			sendMessage("root-size{session=" + sessionID + ",width=" + rect.width + ",height=" + rect.height +"}");
		}
	}
	
	var views = document.getElementsByClassName("ruiView");
	if (views) {
		var message = "resize{session=" + sessionID + ",views=["
		var count = 0
		for (var i = 0; i < views.length; i++) {
			let element = views[i];
			let noresize = element.getAttribute("data-noresize");
			if (!noresize) {
				let rect = element.getBoundingClientRect();
				let top = getIntAttribute(element, "data-top");
				let left = getIntAttribute(element, "data-left");
				let width = getIntAttribute(element, "data-width");
				let height = getIntAttribute(element, "data-height");
				if (rect.width > 0 && rect.height > 0 && 
					(width != rect.width || height != rect.height || left != rect.left || top != rect.top)) {
					element.setAttribute("data-top", rect.top);
					element.setAttribute("data-left", rect.left);
					element.setAttribute("data-width", rect.width);
					element.setAttribute("data-height", rect.height);
					if (count > 0) {
						message += ",";
					}
					message += "view{id=" + element.id + ",x=" + rect.left + ",y=" + rect.top + ",width=" + rect.width + ",height=" + rect.height + 
						",scroll-x=" + element.scrollLeft + ",scroll-y=" + element.scrollTop + ",scroll-width=" + element.scrollWidth + ",scroll-height=" + element.scrollHeight + "}";
					count += 1;
				}
			}
		}

		if (count > 0) {
			sendMessage(message + "]}");
		}
	}
}

function scrollEvent(element, event) {
	sendMessage("scroll{session=" + sessionID + ",id=" + element.id + ",x=" + element.scrollLeft + 
		",y=" + element.scrollTop + ",width=" + element.scrollWidth + ",height=" + element.scrollHeight + "}");
}

function updateCSSRule(selector, ruleText) {
	var styleSheet = document.styleSheets[0];
	var rules = styleSheet.cssRules ? styleSheet.cssRules : styleSheet.rules
	selector = "." + selector
	for (var i = 0; i < rules.length; i++) {
		var rule = rules[i]
		if (!rule.selectorText) {
			continue;
		}
		if (rule.selectorText == selector) {
			if (styleSheet.deleteRule) {
				styleSheet.deleteRule(i)
			} else if (styleSheet.removeRule) {
				styleSheet.removeRule(i)
			}
			break;
		}
	}
	if (styleSheet.insertRule) {
		styleSheet.insertRule(selector + " { " + ruleText + "}")
	} else if (styleSheet.addRule) {
		styleSheet.addRule(selector, ruleText, rules.length)
	}
	scanElementsSize();
}

function updateCSSStyle(elementId, style) {
	var element = document.getElementById(elementId);
	if (element) {
		element.style = style;
		scanElementsSize();
	}
}

function updateCSSProperty(elementId, property, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.style[property] = value;
		scanElementsSize();
	}
}

function updateProperty(elementId, property, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.setAttribute(property, value);
		scanElementsSize();
	}
}

function removeProperty(elementId, property) {
	var element = document.getElementById(elementId);
	if (element && element.hasAttribute(property)) {
		element.removeAttribute(property);
		scanElementsSize();
	}
}

function updateInnerHTML(elementId, content) {
	var element = document.getElementById(elementId);
	if (element) {
		element.innerHTML = content;
		scanElementsSize();
	}
}

function appendToInnerHTML(elementId, content) {
	var element = document.getElementById(elementId);
	if (element) {
		element.innerHTML += content;
		scanElementsSize();
	}
}

function setDisabled(elementId, disabled) {
	var element = document.getElementById(elementId);
	if (element) {
		if ('disabled' in element) {
			element.disabled = disabled
		} else {
			element.setAttribute("data-disabled", disabled ? "1" : "0");
		}
		scanElementsSize();
	}
}

function focusEvent(element, event) {
	event.stopPropagation();
	sendMessage("focus-event{session=" + sessionID + ",id=" + element.id + "}");
}

function blurEvent(element, event) {
	event.stopPropagation();
	sendMessage("lost-focus-event{session=" + sessionID + ",id=" + element.id + "}");
}

function enterOrSpaceKeyClickEvent(event) {
	if (event.key) {
		return (event.key == " " || event.key == "Enter");
	} else if (event.keyCode) {
		return (event.keyCode == 32 || event.keyCode == 13);
	}
	return false;
}

function activateTab(layoutId, tabNumber) {
	var element = document.getElementById(layoutId);
	if (element) {
		var currentNumber = element.getAttribute("data-current");
		if (currentNumber != tabNumber) {
			function setTab(number, styleProperty, display) {
				var tab = document.getElementById(layoutId + '-' + number);
				if (tab) {	
					tab.className = element.getAttribute(styleProperty);
					var page = document.getElementById(tab.getAttribute("data-view"));
					if (page) {
						page.style.display = display;
					}
					return
				}
				var page = document.getElementById(layoutId + "-page" + number);
				if (page) {
					page.style.display = display;
				}
			}
			setTab(currentNumber, "data-inactiveTabStyle", "none")
			setTab(tabNumber, "data-activeTabStyle", "");
			element.setAttribute("data-current", tabNumber);
			scanElementsSize()
		}
	}
}

function tabClickEvent(tab, layoutId, tabNumber, event) {
	event.stopPropagation();
	event.preventDefault();
	activateTab(layoutId, tabNumber)
	if (tab) {
		tab.blur()
	}
	sendMessage("tabClick{session=" + sessionID + ",id=" + layoutId + ",number=" + tabNumber + "}");
}

function tabKeyClickEvent(layoutId, tabNumber, event) {
	if (enterOrSpaceKeyClickEvent(event)) {
		tabClickEvent(null, layoutId, tabNumber, event)
	}
}

function tabCloseClickEvent(button, layoutId, tabNumber, event) {
	event.stopPropagation();
	event.preventDefault();
	if (button) {
		button.blur()
	}
	sendMessage("tabCloseClick{session=" + sessionID + ",id=" + layoutId + ",number=" + tabNumber + "}");
}

function tabCloseKeyClickEvent(layoutId, tabNumber, event) {
	if (enterOrSpaceKeyClickEvent(event)) {
		tabCloseClickEvent(null, layoutId, tabNumber, event)
	}
}


function keyEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id;
	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.key) {
		switch (event.key) {
			case '"':
				message += ",key=`\"`";
				break

			case '\\':
				message += ",key=`\\`";
				break

			default:
				message += ",key=\"" + event.key + "\"";
		}
	}
	if (event.code) {
		message += ",code=\"" + event.code + "\"";
	}
	if (event.repeat) {
		message += ",repeat=1";
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}

	message += "}"
	sendMessage(message);
}

function keyDownEvent(element, event) {
	keyEvent(element, event, "key-down-event")
}

function keyUpEvent(element, event) {
	keyEvent(element, event, "key-up-event")
}

function mouseEventData(element, event) {
	var message = ""

	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.button) {
		message += ",button=" + event.button;
	}
	if (event.buttons) {
		message += ",buttons=" + event.buttons;
	}
	if (event.clientX) {
		var x = event.clientX;
		var el = element;
		if (el.parentElement) {
			x += el.parentElement.scrollLeft;
		}
		while (el) {
			x -= el.offsetLeft
			el = el.parentElement
		}

		message += ",x=" + x + ",clientX=" + event.clientX;
	}
	if (event.clientY) {
		var y = event.clientY;
		var el = element;
		if (el.parentElement) {
			y += el.parentElement.scrollTop;
		}
		while (el) {
			y -= el.offsetTop
			el = el.parentElement
		}

		message += ",y=" + y + ",clientY=" + event.clientY;
	}
	if (event.screenX) {
		message += ",screenX=" + event.screenX;
	}
	if (event.screenY) {
		message += ",screenY=" + event.screenY;
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}
	return message
}

function mouseEvent(element, event, tag) {
	event.stopPropagation();
	//event.preventDefault()

	var message = tag + "{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event) + "}";
	sendMessage(message);
}

function mouseDownEvent(element, event) {
	mouseEvent(element, event, "mouse-down")
}

function mouseUpEvent(element, event) {
	mouseEvent(element, event, "mouse-up")
}

function mouseMoveEvent(element, event) {
	mouseEvent(element, event, "mouse-move")
}

function mouseOverEvent(element, event) {
	mouseEvent(element, event, "mouse-over")
}

function mouseOutEvent(element, event) {
	mouseEvent(element, event, "mouse-out")
}

function clickEvent(element, event) {
	mouseEvent(element, event, "click-event")
	event.preventDefault();
	event.stopPropagation();
}

function doubleClickEvent(element, event) {
	mouseEvent(element, event, "double-click-event")
	event.preventDefault();
}

function contextMenuEvent(element, event) {
	mouseEvent(element, event, "context-menu-event")
	event.preventDefault();
}

function pointerEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event);

	if (event.pointerId) {
		message += ",pointerId=" + event.pointerId;
	}
	if (event.width) {
		message += ",width=" + event.width;
	}
	if (event.height) {
		message += ",height=" + event.height;
	}
	if (event.pressure) {
		message += ",pressure=" + event.pressure;
	}
	if (event.tangentialPressure) {
		message += ",tangentialPressure=" + event.tangentialPressure;
	}
	if (event.tiltX) {
		message += ",tiltX=" + event.tiltX;
	}
	if (event.tiltY) {
		message += ",tiltY=" + event.tiltY;
	}
	if (event.twist) {
		message += ",twist=" + event.twist;
	}
	if (event.pointerType) {
		message += ",pointerType=" + event.pointerType;
	}
	if (event.isPrimary) {
		message += ",isPrimary=1";
	}

	message += "}";
	sendMessage(message);
}

function pointerDownEvent(element, event) {
	pointerEvent(element, event, "pointer-down")
}

function pointerUpEvent(element, event) {
	pointerEvent(element, event, "pointer-up")
}

function pointerMoveEvent(element, event) {
	pointerEvent(element, event, "pointer-move")
}

function pointerCancelEvent(element, event) {
	pointerEvent(element, event, "pointer-cancel")
}

function pointerOverEvent(element, event) {
	pointerEvent(element, event, "pointer-over")
}

function pointerOutEvent(element, event) {
	pointerEvent(element, event, "pointer-out")
}

function touchEvent(element, event, tag) {
	event.stopPropagation();

	var message = tag + "{session=" + sessionID + ",id=" + element.id;
	if (event.timeStamp) {
		message += ",timeStamp=" + event.timeStamp;
	}
	if (event.touches && event.touches.length > 0) {
		message += ",touches=["
		for (var i = 0; i < event.touches.length; i++) {
			var touch = event.touches.item(i)
			if (touch) {
				if (i > 0) {
					message += ","	
				}
				message += "touch{identifier=" + touch.identifier;

				var x = touch.clientX;
				var y = touch.clientY;
				var el = element;
				if (el.parentElement) {
					x += el.parentElement.scrollLeft;
					y += el.parentElement.scrollTop;
				}
				while (el) {
					x -= el.offsetLeft
					y -= el.offsetTop
					el = el.parentElement
				}
		
				message += ",x=" + x + ",y=" + y + ",clientX=" + touch.clientX + ",clientY=" + touch.clientY +
					",screenX=" + touch.screenX + ",screenY=" + touch.screenY + ",radiusX=" + touch.radiusX +
					",radiusY=" + touch.radiusY + ",rotationAngle=" + touch.rotationAngle + ",force=" + touch.force + "}"
			}
		}
		message += "]"
	}
	if (event.ctrlKey) {
		message += ",ctrlKey=1";
	}
	if (event.shiftKey) {
		message += ",shiftKey=1";
	}
	if (event.altKey) {
		message += ",altKey=1";
	}
	if (event.metaKey) {
		message += ",metaKey=1";
	}

	message += "}";
	sendMessage(message);
}

function touchStartEvent(element, event) {
	touchEvent(element, event, "touch-start")
}

function touchEndEvent(element, event) {
	touchEvent(element, event, "touch-end")
}

function touchMoveEvent(element, event) {
	touchEvent(element, event, "touch-move")
}

function touchCancelEvent(element, event) {
	touchEvent(element, event, "touch-cancel")
}

function dropDownListEvent(element, event) {
	event.stopPropagation();
	var message = "itemSelected{session=" + sessionID + ",id=" + element.id + ",number=" + element.selectedIndex.toString() + "}"
	sendMessage(message);
}

function selectDropDownListItem(elementId, number) {
	var element = document.getElementById(elementId);
	if (element) {
		element.selectedIndex = number;
		scanElementsSize();
	}
}

function listItemClickEvent(element, event) {
	event.stopPropagation();
	
	if (element.getAttribute("data-disabled") == "1") {
		return
	}

	var selected = false;
	if (element.classList) {
		const focusStyle = getListFocusedItemStyle(element);
		const blurStyle = getListSelectedItemStyle(element);
		selected = (element.classList.contains(focusStyle) || element.classList.contains(blurStyle));
	} 

	var list = element.parentNode.parentNode
	if (list) {
		if (!selected) {
			selectListItem(list, element, true)
		}

		var message = "itemClick{session=" + sessionID + ",id=" + list.id + "}"
		sendMessage(message);
	}
}

function getListItemNumber(itemId) {
	var pos = itemId.indexOf("-")
	if (pos >= 0) {
		return parseInt(itemId.substring(pos+1))
	}
}

function getStyleAttribute(element, attr, defValue) {
	var result = element.getAttribute(attr);
	if (result) {
		return result;
	}
	return defValue;
}

function getListFocusedItemStyle(element) {
	return getStyleAttribute(element, "data-focusitemstyle", "ruiListItemFocused");
}

function getListSelectedItemStyle(element) {
	return getStyleAttribute(element, "data-bluritemstyle", "ruiListItemSelected");
}

function selectListItem(element, item, needSendMessage) {
	var currentId = element.getAttribute("data-current");
	var message;
	const focusStyle = getListFocusedItemStyle(element);
	const blurStyle = getListSelectedItemStyle(element);

	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			if (current.classList) {
				current.classList.remove(focusStyle, blurStyle);
			}
			if (sendMessage) {
				message = "itemUnselected{session=" + sessionID + ",id=" + element.id + "}";
			}
		}
	}

	if (item) {
		if (element === document.activeElement) {
			if (item.classList) {
				item.classList.add(focusStyle);
			}
		} else {
			if (item.classList) {
				item.classList.add(blurStyle);
			}
		}

		element.setAttribute("data-current", item.id);
		if (sendMessage) {
			var number = getListItemNumber(item.id)
			if (number != undefined) {
				message = "itemSelected{session=" + sessionID + ",id=" + element.id + ",number=" + number + "}";
			}
		}

		if (item.scrollIntoViewIfNeeded) {
			item.scrollIntoViewIfNeeded()
		} else {
			item.scrollIntoView({block: "nearest", inline: "nearest"});
		}
	/*
		var left = item.offsetLeft - element.offsetLeft;
		if (left < element.scrollLeft) {
			element.scrollLeft = left;
		}

		var top = item.offsetTop - element.offsetTop;
		if (top < element.scrollTop) {
			element.scrollTop = top;
		}
		
		var right = left + item.offsetWidth;
		if (right > element.scrollLeft + element.clientWidth) {
			element.scrollLeft = right - element.clientWidth;
		}

		var bottom = top + item.offsetHeight
		if (bottom > element.scrollTop + element.clientHeight) {
			element.scrollTop = bottom - element.clientHeight;
		}*/
	}

	if (needSendMessage && message != undefined) {
		sendMessage(message);
	}
	scanElementsSize();
}

function findRightListItem(list, x, y) {
	list = list.childNodes[0];
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.getAttribute("data-disabled") == "1") {
			continue;
		}
		if (item.offsetLeft >= x) {
			if (result) {
				var result_dy = Math.abs(result.offsetTop - y);
				var item_dy = Math.abs(item.offsetTop - y);
				if (item_dy < result_dy || (item_dy == result_dy && (item.offsetLeft - x) < (result.offsetLeft - x))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findLeftListItem(list, x, y) {
	list = list.childNodes[0];
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.getAttribute("data-disabled") == "1") {
			continue;
		}
		if (item.offsetLeft < x) {
			if (result) {
				var result_dy = Math.abs(result.offsetTop - y);
				var item_dy = Math.abs(item.offsetTop - y);
				if (item_dy < result_dy || (item_dy == result_dy && (x - item.offsetLeft) < (x - result.offsetLeft))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findTopListItem(list, x, y) {
	list = list.childNodes[0];
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.getAttribute("data-disabled") == "1") {
			continue;
		}
		if (item.offsetTop < y) {
			if (result) {
				var result_dx = Math.abs(result.offsetLeft - x);
				var item_dx = Math.abs(item.offsetLeft - x);
				if (item_dx < result_dx || (item_dx == result_dx && (y - item.offsetTop) < (y - result.offsetTop))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function findBottomListItem(list, x, y) {
	list = list.childNodes[0];
	var result;
	var count = list.childNodes.length;
	for (var i = 0; i < count; i++) {
		var item = list.childNodes[i];
		if (item.getAttribute("data-disabled") == "1") {
			continue;
		}
		if (item.offsetTop >= y) {
			if (result) {
				var result_dx = Math.abs(result.offsetLeft - x);
				var item_dx = Math.abs(item.offsetLeft - x);
				if (item_dx < result_dx || (item_dx == result_dx && (item.offsetTop - y) < (result.offsetTop - y))) {
					result = item;	
				}
			} else {
				result = item;
			}
		}
	}
	return result
}

function getKey(event) {
	if (event.key) {
		return event.key;
	}

	if (event.keyCode) {
		switch (event.keyCode) {
			case 13: return "Enter";
			case 32: return " ";
			case 33: return "PageUp";
			case 34: return "PageDown";
			case 35: return "End";
			case 36: return "Home";
			case 37: return "ArrowLeft";
			case 38: return "ArrowUp";
			case 39: return "ArrowRight";
			case 40: return "ArrowDown";
		}
	}
}

function listViewKeyDownEvent(element, event) {
	const key = getKey(event);
	if (key) {
		var currentId = element.getAttribute("data-current");
		var current
		if (currentId) {
			current = document.getElementById(currentId);
			//number = getListItemNumber(currentId);
		}
		if (current) {
			var item
			switch (key) {
				case " ": 
				case "Enter":
					var message = "itemClick{session=" + sessionID + ",id=" + element.id + "}";
					sendMessage(message);
					break;

				case "ArrowLeft":
					item = findLeftListItem(element, current.offsetLeft, current.offsetTop);
					break;
				
				case "ArrowRight":
					item = findRightListItem(element, current.offsetLeft + current.offsetWidth, current.offsetTop);
					break;
		
				case "ArrowDown":
					item = findBottomListItem(element, current.offsetLeft, current.offsetTop + current.offsetHeight);
					break;

				case "ArrowUp":
					item = findTopListItem(element, current.offsetLeft, current.offsetTop);
					break;

				case "Home":
					item = element.childNodes[0];
					break;

				case "End":
					item = element.childNodes[element.childNodes.length - 1];
					break;

				case "PageUp":
					// TODO
					break;

				case "PageDown":
					// TODO
					break;

				default:
					return;
			}
			if (item && item !== current) {
				selectListItem(element, item, true);
			}
		} else {
			switch (key) {
				case " ": 
				case "Enter":
				case "ArrowLeft":
				case "ArrowUp":
				case "ArrowRight":
				case "ArrowDown":
				case "Home":
				case "End":
				case "PageUp":
				case "PageDown":
					var list = element.childNodes[0];
					var count = list.childNodes.length;
					for (var i = 0; i < count; i++) {
						var item = list.childNodes[i];
						if (item.getAttribute("data-disabled") == "1") {
							continue;
						}
						selectListItem(element, item, true);
						return;
					}
					break;

				default:
					return;
			}
		}
	}

	event.stopPropagation();
	event.preventDefault();
}

function listViewFocusEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			if (current.classList) {
				current.classList.remove(getListSelectedItemStyle(element));
				current.classList.add(getListFocusedItemStyle(element));
			} 
		}
	}
}

function listViewBlurEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			if (current.classList) {
				current.classList.remove(getListFocusedItemStyle(element));
				current.classList.add(getListSelectedItemStyle(element));
			}
		}
	}
}

function selectRadioButton(radioButtonId) {
	var element = document.getElementById(radioButtonId);
	if (element) {
		var list = element.parentNode
		if (list) {
			var current = list.getAttribute("data-current");
			if (current) {
				if (current === radioButtonId) {
					return
				}

				var mark = document.getElementById(current + "mark");
				if (mark) {
					//mark.hidden = true
					mark.style.visibility = "hidden"
				}
			}

			var mark = document.getElementById(radioButtonId + "mark");
			if (mark) {
				//mark.hidden = false
				mark.style.visibility = "visible"
			}
			list.setAttribute("data-current", radioButtonId);
			var message = "radioButtonSelected{session=" + sessionID + ",id=" + list.id + ",radioButton=" + radioButtonId + "}"
			sendMessage(message);
			scanElementsSize();
		}
	}
}

function unselectRadioButtons(radioButtonsId) {
	var list = document.getElementById(radioButtonsId);
	if (list) {
		var current = list.getAttribute("data-current");
		if (current) {
			var mark = document.getElementById(current + "mark");
			if (mark) {
				mark.style.visibility = "hidden"
			}

			list.removeAttribute("data-current");
		}

		var message = "radioButtonUnselected{session=" + sessionID + ",id=" + list.id + "}"
		sendMessage(message);
		scanElementsSize();
	}
}

function radioButtonClickEvent(element, event) {
	event.stopPropagation();
	event.preventDefault();
	selectRadioButton(element.id)
}

function radioButtonKeyClickEvent(element, event) {
	if (enterOrSpaceKeyClickEvent(event)) {
		radioButtonClickEvent(element, event);
	}
}

function editViewInputEvent(element) {
	var text = element.value
	text = text.replaceAll(/\\/g, "\\\\")
	text = text.replaceAll(/\"/g, "\\\"")
	var message = "textChanged{session=" + sessionID + ",id=" + element.id + ",text=\"" + text + "\"}"
	sendMessage(message);
}

function setInputValue(elementId, text) {
	var element = document.getElementById(elementId);
	if (element) {
		element.value = text;
		scanElementsSize();
	}
}

function fileSelectedEvent(element) {
	var files = element.files;
	if (files) {
		var message = "fileSelected{session=" + sessionID + ",id=" + element.id + ",files=[";
		for(var i = 0; i < files.length; i++) {
			if (i > 0) {
				message += ",";
			}
			message += "_{name=\"" + files[i].name + 
				"\",last-modified=" + files[i].lastModified +
				",size=" + files[i].size +
				",mime-type=\"" + files[i].type + "\"}";
		}
		sendMessage(message + "]}");
	}
}

function loadSelectedFile(elementId, index) {
	var element = document.getElementById(elementId);
	if (element) {
		var files = element.files;
		if (files && index >= 0 && index < files.length) {
			const reader = new FileReader();
         	reader.onload = function() { 
				sendMessage("fileLoaded{session=" + sessionID + ",id=" + element.id + 
					",index=" + index + 
					",name=\"" + files[index].name + 
					"\",last-modified=" + files[index].lastModified +
					",size=" + files[index].size +
					",mime-type=\"" + files[index].type + 
					"\",data=`" + reader.result + "`}");
			}
         	reader.onerror = function(error) {
				sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`" + error + "`}");
			}
			reader.readAsDataURL(files[index]);
		} else {
			sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`File not found`}");
		}
	} else {
		sendMessage("fileLoadingError{session=" + sessionID + ",id=" + element.id + ",index=" + index + ",error=`Invalid FilePicker id`}");
	}
}

function startResize(element, mx, my, event) {
	var view = element.parentNode;
	if (!view) {
		return;
	}

	var startX = event.clientX;
	var startY = event.clientY;
	var startWidth = view.offsetWidth
	var startHeight = view.offsetHeight

	document.addEventListener("mousemove", moveHandler, true);
	document.addEventListener("mouseup", upHandler, true);
	
	event.stopPropagation();
	event.preventDefault();

	function moveHandler(e) {
		if (mx != 0) {
			var width = startWidth + (e.clientX - startX) * mx;
			if (width <= 0) {
				width = 1;
			}
			view.style.width = width + "px";
			sendMessage("widthChanged{session=" + sessionID + ",id=" + view.id + ",width=" + view.style.width + "}");
		}
		
		if (my != 0) {
			var height = startHeight + (e.clientY - startY) * my;
			if (height <= 0) {
				height = 1;
			}
			view.style.height = height + "px";
			sendMessage("heightChanged{session=" + sessionID + ",id=" + view.id + ",height=" + view.style.height + "}");
		}

		event.stopPropagation();
		event.preventDefault();
		scanElementsSize();
	}

	function upHandler (e) {
		document.removeEventListener("mouseup", upHandler, true);
		document.removeEventListener("mousemove", moveHandler, true);
		e.stopPropagation();
	}
}

function transitionStartEvent(element, event) {
	var message = "transition-start-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionRunEvent(element, event) {
	var message = "transition-run-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionEndEvent(element, event) {
	var message = "transition-end-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function transitionCancelEvent(element, event) {
	var message = "transition-cancel-event{session=" + sessionID + ",id=" + element.id; 
	if (event.propertyName) {
		message += ",property=" + event.propertyName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationStartEvent(element, event) {
	var message = "animation-start-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationEndEvent(element, event) {
	var message = "animation-end-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationCancelEvent(element, event) {
	var message = "animation-cancel-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function animationIterationEvent(element, event) {
	var message = "animation-iteration-event{session=" + sessionID + ",id=" + element.id; 
	if (event.animationName) {
		message += ",name=" + event.animationName
	}
	sendMessage(message + "}");
	event.stopPropagation();
}

function stackTransitionEndEvent(stackId, propertyName, event) {
	sendMessage("transition-end-event{session=" + sessionID + ",id=" + stackId + ",property=" + propertyName + "}");
	event.stopPropagation();
}

var images = new Map();

function loadImage(url) {
	var img = images.get(url);
	if (img != undefined) {
		return
	}
	
	img = new Image(); 
	img.addEventListener("load", function() {
		images.set(url, img)
		var message = "imageLoaded{session=" + sessionID + ",url=\"" + url + "\""; 
		if (img.naturalWidth) {
			message += ",width=" + img.naturalWidth
		}
		if (img.naturalHeight) {
			message += ",height=" + img.naturalHeight
		}
		sendMessage(message + "}")
	}, false);

	img.addEventListener("error", function(event) {
		var message = "imageError{session=" + sessionID + ",url=\"" + url + "\""; 
		if (event && event.message) {
			var text = event.message.replaceAll(new RegExp("\"", 'g'), "\\\"")
			message += ",message=\"" + text + "\""; 
		}
		sendMessage(message + "}")
	}, false);

	img.src = url;
}

function loadInlineImage(url, content) {
	var img = images.get(url);
	if (img != undefined) {
		return
	}
	
	img = new Image(); 
	img.addEventListener("load", function() {
		images.set(url, img)
		var message = "imageLoaded{session=" + sessionID + ",url=\"" + url + "\""; 
		if (img.naturalWidth) {
			message += ",width=" + img.naturalWidth
		}
		if (img.naturalHeight) {
			message += ",height=" + img.naturalHeight
		}
		sendMessage(message + "}")
	}, false);

	img.addEventListener("error", function(event) {
		var message = "imageError{session=" + sessionID + ",url=\"" + url + "\""; 
		if (event && event.message) {
			var text = event.message.replaceAll(new RegExp("\"", 'g'), "\\\"")
			message += ",message=\"" + text + "\""; 
		}
		sendMessage(message + "}")
	}, false);

	img.src = content;
}

function clickClosePopup(element, e) {
	var popupId = element.getAttribute("data-popupId");
	sendMessage("clickClosePopup{session=" + sessionID + ",id=" + popupId + "}")
	e.stopPropagation();
}

function scrollTo(elementId, x, y) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(x, y);
	}
}

function scrollToStart(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(0, 0);
	}
}

function scrollToEnd(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.scrollTo(0, element.scrollHeight - element.offsetHeight);
	}
}

function focus(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.focus();
	}
}

function blur(elementId) {
	var element = document.getElementById(elementId);
	if (element) {
		element.blur();
	}
}

function blurCurrent() {
	if (document.activeElement != document.body) {
		document.activeElement.blur();
	}
}

function playerEvent(element, tag) {
	//event.stopPropagation();
	sendMessage(tag + "{session=" + sessionID + ",id=" + element.id + "}");
}

function playerTimeUpdatedEvent(element) {
	var message = "time-update-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.currentTime) {
		message += element.currentTime;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerDurationChangedEvent(element) {
	var message = "duration-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.duration) {
		message += element.duration;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerVolumeChangedEvent(element) {
	var message = "volume-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.volume && !element.muted) {
		message += element.volume;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerRateChangedEvent(element) {
	var message = "rate-changed-event{session=" + sessionID + ",id=" + element.id + ",value=";
	if (element.playbackRate) {
		message += element.playbackRate;
	} else {
		message += "0";
	}
	sendMessage(message + "}");
}

function playerErrorEvent(element) {
	var message = "player-error-event{session=" + sessionID + ",id=" + element.id;
	if (element.error) {
		if (element.error.code) {
			message += ",code=" + element.error.code;
		}
		if (element.error.message) {
			message += ",message=`" + element.error.message + "`";
		}
	}
	sendMessage(message + "}");
}

function setMediaMuted(elementId, value) {
	var element = document.getElementById(elementId);
	if (element) {
		element.muted = value
	}
}

function mediaPlay(elementId) {
	var element = document.getElementById(elementId);
	if (element && element.play) {
		element.play()
	}
}

function mediaPause(elementId) {
	var element = document.getElementById(elementId);
	if (element && element.pause) {
		element.pause()
	}
}

function mediaSetSetCurrentTime(elementId, time) {
	var element = document.getElementById(elementId);
	if (element) {
		element.currentTime = time
	}
}

function mediaSetPlaybackRate(elementId, time) {
	var element = document.getElementById(elementId);
	if (element) {
		element.playbackRate = time
	}
}

function mediaSetVolume(elementId, volume) {
	var element = document.getElementById(elementId);
	if (element) {
		element.volume = volume
	}
}

function startDownload(url, filename) {
	var element = document.getElementById("ruiDownloader");
	if (element) {
		element.href = url;
		element.setAttribute("download", filename);
		element.click();
	}
}

function setTitle(title) {
	document.title = title;
}

function setTitleColor(color) {
	var metas = document.getElementsByTagName('meta');
	if (metas) {
		var item = metas.namedItem('theme-color');
		if (item) {
			item.setAttribute('content', color)
			return
		}
	}

	var meta = document.createElement('meta');
	meta.setAttribute('name', 'theme-color');
	meta.setAttribute('content', color);
	document.getElementsByTagName('head')[0].appendChild(meta);
}

function openURL(url) {
	window.open(url, "_blank");
}

function detailsEvent(element) {
	sendMessage("details-open{session=" + sessionID + ",id=" + element.id + ",open=" + (element.open ? "1}" : "0}"));
}

function getTableFocusedItemStyle(element) {
	return getStyleAttribute(element, "data-focusitemstyle", "ruiCurrentTableCellFocused");
}

function getTableSelectedItemStyle(element) {
	return getStyleAttribute(element, "data-bluritemstyle", "ruiCurrentTableCell");
}

function tableViewFocusEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current) {
			if (current.classList) {
				current.classList.remove(getTableSelectedItemStyle(element));
				current.classList.add(getTableFocusedItemStyle(element));
			} 
		}
	}
}

function tableViewBlurEvent(element, event) {
	var currentId = element.getAttribute("data-current");
	if (currentId) {
		var current = document.getElementById(currentId);
		if (current && current.classList) {
			current.classList.remove(getTableFocusedItemStyle(element));
			current.classList.add(getTableSelectedItemStyle(element));
		}
	}
}

function setTableCellCursorByID(tableID, row, column) {
	var table = document.getElementById(tableID);
	if (table) {
		if (!setTableCellCursor(table, row, column)) {
			const focusStyle = getTableFocusedItemStyle(table);
			const oldCellID = table.getAttribute("data-current");
			if (oldCellID) {
				const oldCell = document.getElementById(oldCellID);
				if (oldCell && oldCell.classList) {
					oldCell.classList.remove(focusStyle);
					oldCell.classList.remove(getTableSelectedItemStyle(table));
				}
				table.removeAttribute("data-current");
			}
		}
	}
}

function setTableCellCursor(element, row, column) {
	const cellID = element.id + "-" + row + "-" + column;
	var cell = document.getElementById(cellID);
	if (!cell || cell.getAttribute("data-disabled")) {
		return false;
	}

	const focusStyle = getTableFocusedItemStyle(element);
	const oldCellID = element.getAttribute("data-current");
	if (oldCellID) {
		const oldCell = document.getElementById(oldCellID);
		if (oldCell && oldCell.classList) {
			oldCell.classList.remove(focusStyle);
			oldCell.classList.remove(getTableSelectedItemStyle(element));
		}
	}

	cell.classList.add(focusStyle);
	element.setAttribute("data-current", cellID);
	if (cell.scrollIntoViewIfNeeded) {
		cell.scrollIntoViewIfNeeded()
	} else {
		cell.scrollIntoView({block: "nearest", inline: "nearest"});
	}

	sendMessage("currentCell{session=" + sessionID + ",id=" + element.id + 
			",row=" + row + ",column=" + column + "}");
	return true;
}

function moveTableCellCursor(element, row, column, dr, dc) {
	const rows = element.getAttribute("data-rows");
	if (!rows) {
		return;
	}
	const columns = element.getAttribute("data-columns");
	if (!columns) {
		return;
	}

	const rowCount = parseInt(rows);
	const columnCount = parseInt(columns);
	
	row += dr;
	column += dc;
	while (row >= 0 && row < rowCount && column >= 0 && column < columnCount) {
		if (setTableCellCursor(element, row, column)) {
			return;
		} else if (dr == 0) {
			var r2 = row - 1;
			while (r2 >= 0) {
				if (setTableCellCursor(element, r2, column)) {
					return;
				}
				r2--;
			}
		} else if (dc == 0) {
			var c2 = column - 1;
			while (c2 >= 0) {
				if (setTableCellCursor(element, row, c2)) {
					return;
				}
				c2--;
			}
		}
		row += dr;
		column += dc;
	}
}

function tableViewCellKeyDownEvent(element, event) {
	const key = getKey(event);
	if (!key) {
		return;
	}

	const currentId = element.getAttribute("data-current");
	if (!currentId || currentId == "") {
		switch (key) {
			case "ArrowLeft":
			case "ArrowRight":
			case "ArrowDown":
			case "ArrowUp":
			case "Home":
			case "End":
			case "PageUp":
			case "PageDown":
				event.stopPropagation();
				event.preventDefault();
				const rows = element.getAttribute("data-rows");
				const columns = element.getAttribute("data-columns");
				if (rows && columns) {
					const rowCount = parseInt(rows);
					const columnCount = parseInt(rows);
					row = 0;
					while (row < rowCount) {
						column = 0;
						while (columns < columnCount) {
							if (setTableCellCursor(element, row, column)) {
								return;
							}
							column++;
						}
						row++;
					}
				}
				break;
		}
		return;
	}

	const elements = currentId.split("-");
	if (elements.length >= 3) {
		const row = parseInt(elements[1], 10)
		const column = parseInt(elements[2], 10)

		switch (key) {
			case " ": 
			case "Enter":
				sendMessage("cellClick{session=" + sessionID + ",id=" + element.id + 
							",row=" + row + ",column=" + column + "}");
				break;

			case "ArrowLeft":
				moveTableCellCursor(element, row, column, 0, -1)
				break;
			
			case "ArrowRight":
				moveTableCellCursor(element, row, column, 0, 1)
				break;

			case "ArrowDown":
				moveTableCellCursor(element, row, column, 1, 0)
				break;

			case "ArrowUp":
				moveTableCellCursor(element, row, column, -1, 0)
				break;

			case "Home":
				// TODO
				break;

			case "End":
				/*var newRow = rowCount-1;
				while (newRow > row) {
					if (setTableRowCursor(element, newRow)) {
						break;
					}
					newRow--;
				}*/
				// TODO
				break;

			case "PageUp":
				// TODO
				break;

			case "PageDown":
				// TODO
				break;

			default:
				return;
		}

		event.stopPropagation();
		event.preventDefault();
	} else {
		element.setAttribute("data-current", "");
	}
}

function setTableRowCursorByID(tableID, row) {
	var table = document.getElementById(tableID);
	if (table) {
		if (!setTableRowCursor(table, row)) {
			const focusStyle = getTableFocusedItemStyle(table);
			const oldRowID = table.getAttribute("data-current");
			if (oldRowID) {
				const oldRow = document.getElementById(oldRowID);
				if (oldRow && oldRow.classList) {
					oldRow.classList.remove(focusStyle);
					oldRow.classList.remove(getTableSelectedItemStyle(table));
				}
				table.removeAttribute("data-current");
			}
		}
	}
}

function setTableRowCursor(element, row) {
	const tableRowID = element.id + "-" + row;
	var tableRow = document.getElementById(tableRowID);
	if (!tableRow || tableRow.getAttribute("data-disabled")) {
		return false;
	}

	const focusStyle = getTableFocusedItemStyle(element);
	const oldRowID = element.getAttribute("data-current");
	if (oldRowID) {
		const oldRow = document.getElementById(oldRowID);
		if (oldRow && oldRow.classList) {
			oldRow.classList.remove(focusStyle);
			oldRow.classList.remove(getTableSelectedItemStyle(element));
		}
	}

	tableRow.classList.add(focusStyle);
	element.setAttribute("data-current", tableRowID);
	if (tableRow.scrollIntoViewIfNeeded) {
		tableRow.scrollIntoViewIfNeeded()
	} else {
		tableRow.scrollIntoView({block: "nearest", inline: "nearest"});
	}

	sendMessage("currentRow{session=" + sessionID + ",id=" + element.id + ",row=" + row + "}");
	return true;
}

function moveTableRowCursor(element, row, dr) {
	const rows = element.getAttribute("data-rows");
	if (!rows) {
		return;
	}

	const rowCount = parseInt(rows);
	row += dr;
	while (row >= 0 && row < rowCount) {
		if (setTableRowCursor(element, row)) {
			return;
		}
		row += dr;
	}
}

function tableViewRowKeyDownEvent(element, event) {
	const key = getKey(event);
	if (key) {
		const currentId = element.getAttribute("data-current");
		if (currentId) {
			const elements = currentId.split("-");
			if (elements.length >= 2) {
				const row = parseInt(elements[1], 10);
				switch (key) {
					case " ": 
					case "Enter":
						sendMessage("rowClick{session=" + sessionID + ",id=" + element.id + ",row=" + row + "}");
						break;
		
					case "ArrowDown":
						moveTableRowCursor(element, row, 1)
						break;
		
					case "ArrowUp":
						moveTableRowCursor(element, row, -1)
						break;
		
					case "Home":
						var newRow = 0;
						while (newRow < row) {
							if (setTableRowCursor(element, newRow)) {
								break;
							}
							newRow++;
						}
						break;
		
					case "End":
						var newRow = rowCount-1;
						while (newRow > row) {
							if (setTableRowCursor(element, newRow)) {
								break;
							}
							newRow--;
						}
						break;
		
					case "PageUp":
						// TODO
						break;
		
					case "PageDown":
						// TODO
						break;
		
					default:
						return;
				}
				event.stopPropagation();
				event.preventDefault();
				return;
			}
		}

		switch (key) {
			case "ArrowLeft":
			case "ArrowRight":
			case "ArrowDown":
			case "ArrowUp":
			case "Home":
			case "End":
			case "PageUp":
			case "PageDown":
				const rows = element.getAttribute("data-rows");
				if (rows) {
					const rowCount = parseInt(rows);
					row = 0;
					while (row < rowCount) {
						if (setTableRowCursor(element, row)) {
							break;
						}
						row++;
					}
				}
				break;

			default:
				return;
		}
	} 

	event.stopPropagation();
	event.preventDefault();
}

function tableCellClickEvent(element, event) {
	event.preventDefault();

	const elements = element.id.split("-");
	if (elements.length < 3) {
		return
	}

	const tableID = elements[0];
	const row = parseInt(elements[1], 10);
	const column = parseInt(elements[2], 10);
	const table = document.getElementById(tableID);
	if (table) {
		const selection = table.getAttribute("data-selection");
		if (selection == "cell") {
			const currentID = table.getAttribute("data-current");
			if (!currentID || currentID != element.ID) {
				setTableCellCursor(table, row, column)
			}
		}
	}

	sendMessage("cellClick{session=" + sessionID + ",id=" + tableID + 
				",row=" + row + ",column=" + column + "}");
}

function tableRowClickEvent(element, event) {
	event.preventDefault();

	const elements = element.id.split("-");
	if (elements.length < 2) {
		return
	}

	const tableID = elements[0];
	const row = parseInt(elements[1], 10);
	const table = document.getElementById(tableID);
	if (table) {
		const selection = table.getAttribute("data-selection");
		if (selection == "row") {
			const currentID = table.getAttribute("data-current");
			if (!currentID || currentID != element.ID) {
				setTableRowCursor(table, row)
			}
		}
	}

	sendMessage("rowClick{session=" + sessionID + ",id=" + tableID + ",row=" + row + "}");
}

function imageLoaded(element, event) {
	var message = "imageViewLoaded{session=" + sessionID + ",id=" + element.id +
		",natural-width=" + element.naturalWidth +
		",natural-height=" + element.naturalHeight +
		",current-src=\"" + element.currentSrc +  "\"}";
	sendMessage(message);
	scanElementsSize()
}

function imageError(element, event) {
	var message = "imageViewError{session=" + sessionID + ",id=" + element.id + "}";
	sendMessage(message);
}

function canvasTextMetrics(answerID, elementId, font, text) {
	var w = 0;
	var ascent = 0;
	var descent = 0;
	var left = 0;
	var right = 0;

	const canvas = document.getElementById(elementId);
	if (canvas) {
		const ctx = canvas.getContext('2d');
  		if (ctx) {
			ctx.save()
			const dpr = window.devicePixelRatio || 1;
			ctx.scale(dpr, dpr);
			ctx.font = font;
			ctx.textBaseline = 'alphabetic';
			ctx.textAlign = 'start';
			var metrics = ctx.measureText(text)
			w = metrics.width;
			ascent = metrics.actualBoundingBoxAscent;
			descent = metrics.actualBoundingBoxDescent;
			left = metrics.actualBoundingBoxLeft;
			right = metrics.actualBoundingBoxRight;
			ctx.restore();
		}
	}

	sendMessage('answer{answerID=' + answerID + ', width=' + w + ', ascent=' + ascent + 
		', descent=' + descent + ', left=' + left + ', right=' + right + '}');
}

function getPropertyValue(answerID, elementId, name) {
	const element = document.getElementById(elementId);
	if (element && element[name]) {
		sendMessage('answer{answerID=' + answerID + ', value="' + element[name] + '"}')
		return
	}

	sendMessage('answer{answerID=' + answerID + ', value=""}')
}

function appendStyles(styles) {
	document.querySelector('style').textContent += styles
}

function getCanvasContext(elementId) {
	const canvas = document.getElementById(elementId)
	if (canvas) {
		const ctx = canvas.getContext('2d');
		if (ctx) {
			const dpr = window.devicePixelRatio || 1;
			ctx.canvas.width = dpr * canvas.clientWidth;
			ctx.canvas.height = dpr * canvas.clientHeight;
			ctx.scale(dpr, dpr);
			return ctx;
		}
	}
	return null;
}

function localStorageSet(key, value) {
	try {
		localStorage.setItem(key, value)
	} catch (err) {
		sendMessage("storageError{session=" + sessionID + ", error=`" + err + "`}")
	}
}

function localStorageClear() {
	try {
		localStorage.setItem(key, value)
	} catch (err) {
		sendMessage("storageError{session=" + sessionID + ", error=`" + err + "`}")
	}
}

function showTooltip(element, tooltip) {
	const layer = document.getElementById("ruiTooltipLayer");
	if (!layer) {
		return;
	}

	layer.style.left = "0px";
	layer.style.right = "0px";

	var tooltipBox = document.getElementById("ruiTooltipText");
	if (tooltipBox) {
		tooltipBox.innerHTML = tooltip;
	}

	var left = element.offsetLeft;
	var top = element.offsetTop;
	var width = element.offsetWidth;
	var height = element.offsetHeight;
	var parent = element.offsetParent;

	while (parent) {
		left += parent.offsetLeft;
		top += parent.offsetTop;
		width = parent.offsetWidth;
		height = parent.offsetHeight;
		parent = parent.offsetParent;
	}

	if (element.offsetWidth >= tooltipBox.offsetWidth) {
		layer.style.left = left + "px";
		layer.style.justifyItems = "start";
	} else {
		const rightOff = width - (left + element.offsetWidth);
		if (left > rightOff) {
			if (width - left < tooltipBox.offsetWidth) {
				layer.style.right = rightOff + "px";
				layer.style.justifyItems = "end";
			} else {
				layer.style.left = (left - rightOff) + "px";
				layer.style.right = "0px";
				layer.style.justifyItems = "center";
			}
		} else {
			if (width - rightOff < tooltipBox.offsetWidth) {
				layer.style.left = left + "px";
				layer.style.justifyItems = "start";
			} else {
				layer.style.right = (rightOff - left) + "px";
				layer.style.justifyItems = "center";
			}
		}
	}

	const bottomOff = height - (top + element.offsetHeight);
	var arrow = document.getElementById("ruiTooltipTopArrow");

	if (bottomOff < arrow.offsetHeight + tooltipBox.offsetHeight) {
		if (arrow) {
			arrow.style.visibility = "hidden";
		}

		arrow = document.getElementById("ruiTooltipBottomArrow");
		if (arrow) {
			arrow.style.visibility = "visible";
		}

		layer.style.top = "0px";
		layer.style.bottom = height - top - arrow.offsetHeight / 2 + "px"; 
		layer.style.gridTemplateRows = "1fr auto auto"

	} else {
		if (arrow) {
			arrow.style.visibility = "visible";
		}

		layer.style.top = top + element.offsetHeight - arrow.offsetHeight / 2 + "px";
		layer.style.bottom = "0px";
		layer.style.gridTemplateRows = "auto auto 1fr"

		arrow = document.getElementById("ruiTooltipBottomArrow");
		if (arrow) {
			arrow.style.visibility = "hidden";
		}
	}
	
	layer.style.visibility = "visible";
	layer.style.opacity = 1;
}

function mouseEnterEvent(element, event) {
	event.stopPropagation();

	let tooltip = element.getAttribute("data-tooltip");
	if (tooltip) {
		showTooltip(element, tooltip);
	}

	sendMessage("mouse-enter{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event) + "}");
}

function mouseLeaveEvent(element, event) {
	event.stopPropagation();

	if (element.getAttribute("data-tooltip")) {
		const layer = document.getElementById("ruiTooltipLayer");
		if (layer) {
			layer.style.opacity = 0;
			layer.style.visibility = "hidden";
		}
	}

	sendMessage("mouse-leave{session=" + sessionID + ",id=" + element.id + mouseEventData(element, event) + "}");
}

function hideTooltip() {
	const layer = document.getElementById("ruiTooltipLayer");
	if (layer) {
		layer.style.opacity = 0;
		layer.style.visibility = "hidden";
	}
}

function stopEventPropagation(element, event) {
	event.stopPropagation();
}

function setCssVar(tag, value) {
	const root = document.querySelector(':root');
	if (root) {
		root.style.setProperty(tag, value);
	}
}
