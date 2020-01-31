!function() {
	var shadowHost = document.getElementById('mmt').parentNode;
	var shadowRoot;
	if (shadowHost.attachShadow) {
		shadowRoot = shadowHost.attachShadow({'mode': 'closed'});
	} else {
		shadowHost.dataset.mmt = 'host';
		shadowRoot = shadowHost;
	}
	shadowRoot.innerHTML = '<link rel="stylesheet" href="css/mmt.css"><div id="mmt" spellcheck="false"><div id="mmtChatPanel"><div class="panel"><div class="panelHead"><img id="mmtShowRoomListPanelBtn" src="http://192.168.20.160/static/v2_pvwc65.png"> <div class="textContent"><span id="mmtRoomName"></span> (<span id="mmtRoomUserCount"></span> members) </div></div><div class="panelBody"><div id="mmtMessageOutput"></div></div><div class="panelFoot"><div class="userMessageBox"><div class="userInfo"><img id="mmtUserIcon" class="userIcon"><div id="mmtNickName" class="nickName"></div></div><div class="myMessage"><textarea id="mmtMessageInput"></textarea><button id="mmtSendBtn">Send</button></div></div></div></div></div><div id="mmtRoomListPanel" hidden><div class="panel"><div class="panelHead"><h2 class="textContent">喵喵通</h2></div><div class="panelBody"><ul id="mmtRoomList"></ul><div id="mmtAddRoom"><input id="mmtRoomNameInput" type="text" placeholder="建立頻道"><button id="mmtJoinBtn">+</button><div></div></div></div></div>';
	function onLoad() {
		if (!shadowRoot.querySelector('#mmt')) {
			requestAnimationFrame(onLoad);
			return;
		}
		var mmtChatPanel            = shadowRoot.querySelector('#mmtChatPanel');
		var mmtShowRoomListPanelBtn = shadowRoot.querySelector('#mmtShowRoomListPanelBtn');
		var mmtRoomName             = shadowRoot.querySelector('#mmtRoomName');
		var mmtRoomUserCount        = shadowRoot.querySelector('#mmtRoomUserCount');
		var mmtMessageOutput        = shadowRoot.querySelector('#mmtMessageOutput');
		var mmtUserIcon             = shadowRoot.querySelector('#mmtUserIcon');
		var mmtNickName             = shadowRoot.querySelector('#mmtNickName');
		var mmtMessageInput         = shadowRoot.querySelector('#mmtMessageInput');
		var mmtSendBtn              = shadowRoot.querySelector('#mmtSendBtn');
		var mmtRoomListPanel        = shadowRoot.querySelector('#mmtRoomListPanel');
		var mmtRoomList             = shadowRoot.querySelector('#mmtRoomList');
		var mmtRoomNameInput        = shadowRoot.querySelector('#mmtRoomNameInput');
		var mmtJoinBtn              = shadowRoot.querySelector('#mmtJoinBtn');
		// var scheme = document.location.protocol == 'https:' ? 'wss': 'ws';
		// var port = document.location.port ? ':' + document.location.port: '';
		// var wsURL = scheme + '://' + document.location.hostname + port + '/echo';
		var wsURL = 'ws://35.187.152.214:8080/echo';
		var ws;
		var nickName = '';
		var roomName = '恆星足球';
		var roomMsgList = {};
		function connect() {
			if (ws) {
				return false;
			}
			ws = new WebSocket(wsURL);
			ws.onopen = function(evt) {
				handleNamespaceConnectedConn();
			}
			ws.onclose = function(evt) {
				// console.log('CLOSE');
				ws = null;
			}
			ws.onmessage = function(evt) {
				receivePacketHandle(evt.data);
			}
			ws.onerror = function(evt) {
				// console.log('ERROR: ' + evt.data);
			}
			return true;
		}
		function disconnect() {
			if (ws) {
				var playerLogout = {'cmd': '4', 'idem': Date.now().toString() ,'payload': {'roomName': roomName}};
				ws.send(JSON.stringify(playerLogout));
			}
		}
		function joinRoom(targetRoomName) {
			if (ws && targetRoomName != '') {
				mmtRoomNameInput.value = '';
				var playerEnterRoom = {'cmd': '10', 'idem': targetRoomName ,'payload': {'roomName': targetRoomName}};
				ws.send(JSON.stringify(playerEnterRoom));
			}
		}
		function exitRoom() {
			if (ws) {
				var playerExitRoom = {'cmd': '8', 'idem': roomName ,'payload': {'roomName': roomName}};
				ws.send(JSON.stringify(playerExitRoom));
			}
		}
		function roomSearch() {
			if (ws) {
				var roomSearchCmd = {'cmd': '16', 'idem': Date.now().toString()};
				ws.send(JSON.stringify(roomSearchCmd));
			}
		}
		function userSearch() {
			if (ws) {
				var userSearchCmd = {'cmd': '6', 'idem': roomName ,'payload': {'roomName': roomName}};
				ws.send(JSON.stringify(userSearchCmd));
			}
		}
		function sendMessage() {
			var msg = mmtMessageInput.value;
			if (msg != '') {
				var sendMessageCmd = {'cmd': '80', 'idem': Date.now().toString(), 'payload': {'message': msg, 'style': '1', 'roomInfo': {'roomName': roomName}}};
				mmtMessageInput.value = '';
				ws.send(JSON.stringify(sendMessageCmd));
			}
		}
		var isRoomListLoading = false;
		function openRoomListPanel() {
			if (!isRoomListLoading) {
				if (mmtRoomListPanel.hidden) {
					isRoomListLoading = true;
					roomSearch();
				} else {
					mmtRoomListPanel.hidden = true;
				}
			}
		}
		function addRoomList(roomList) {
			mmtRoomList.innerHTML = '';
			isRoomListLoading = false;
			roomList.sort(function(x,y){return x.roomName>y.roomName?1:-1;});
			for (var i = 0; i < roomList.length; i++) {
				var targetRoomName = roomList[i].roomName;
				var li = document.createElement('li');
				li.innerHTML = '<div class="textContent">' + htmlSafeText(targetRoomName) + '</div>';
				li.id = 'room-'+ targetRoomName;
				li.onclick = clickRoomList;
				mmtRoomList.appendChild(li);
			}
		}
		var htmlSafeText = function(str) {
			this.textContent = str;
			return this.innerHTML;
		}.bind(document.createElement('div'));
		function addSystemMessage(targetRoomName, msg) {
			if (roomMsgList[targetRoomName]) {
				var msgHtml = '<div class="systemMessage">' + htmlSafeText(msg) + '</div>';
				roomMsgList[targetRoomName].html += msgHtml;
				if (roomName == targetRoomName) {
					var chatPanelBody = shadowRoot.querySelector('#mmtChatPanel .panelBody');
					var isAutoScroll = chatPanelBody.scrollHeight - chatPanelBody.clientHeight - chatPanelBody.scrollTop < 1;
					mmtMessageOutput.innerHTML += msgHtml;
					if (isAutoScroll) {
						chatPanelBody.scrollTop = chatPanelBody.scrollHeight;
					}
				}
			}
		}
		function addUserMessage(msgObj, isHistory) {
			var targetRoomName = msgObj.roomInfo.roomName;
			if (roomMsgList[targetRoomName]) {
				var targetNickName = msgObj.from.nickName;
				var targetIcon = msgObj.from.icon;
				var msg = msgObj.message;
				var msgHtml = '';
				msgHtml += '<div class="userMessageBox"><div class="userInfo">';
				msgHtml += nickName == targetNickName ? '</div><div class="myMessage">': '<img class="userIcon" src="' + targetIcon + '"></div><div class="userMessage">';
				msgHtml += '<div class="nickName">' + targetNickName + '</div><div class="bubble">' + htmlSafeText(msg) + '</div></div></div>';
				if (isHistory) {
					msgHtml += roomMsgList[targetRoomName].html;
					roomMsgList[targetRoomName].html = msgHtml;
				} else {
					roomMsgList[targetRoomName].html += msgHtml;
					msgHtml = roomMsgList[targetRoomName].html;
				}
				if (roomName == targetRoomName) {
					var chatPanelBody = shadowRoot.querySelector('#mmtChatPanel .panelBody');
					var isAutoScroll = chatPanelBody.scrollHeight - chatPanelBody.clientHeight - chatPanelBody.scrollTop < 1;
					mmtMessageOutput.innerHTML = msgHtml;
					if (isAutoScroll) {
						chatPanelBody.scrollTop = chatPanelBody.scrollHeight;
					}
				}
			}
		}
		function addRoom(targetRoomName) {
			roomMsgList[targetRoomName] = {'historyUid': '', 'html': ''};
			changeRoom(targetRoomName);
			var chatHistory = {'cmd': '18', 'idem': targetRoomName, 'payload': {'room': {'roomName': targetRoomName}, 'historyUid': ''}};
			ws.send(JSON.stringify(chatHistory));
		}
		function changeRoom(targetRoomName) {
			mmtRoomListPanel.hidden = true;
			if (roomName != targetRoomName) {
				roomName = targetRoomName;
				mmtRoomName.textContent = roomName;
				mmtMessageOutput.innerHTML = roomMsgList[roomName].html;
				userSearch();
			}
		}
		function clickRoomList() {
			var targetRoomName = this.textContent;
			if (roomMsgList[targetRoomName]) {
				changeRoom(targetRoomName)
			} else {
				joinRoom(targetRoomName);
			}
		}
		function clickJoinRoom() {
			var targetRoomName = mmtRoomNameInput.value;
			if (roomMsgList[targetRoomName]) {
				changeRoom(targetRoomName)
			} else {
				joinRoom(targetRoomName);
			}
		}
		function handleError(reason) {
			// console.log(reason);
			window.alert('error: see the dev console');
		}
		function handleNamespaceConnectedConn() {
			nickName = 'guest' + Math.random().toFixed(4).slice(2);
			var loginInfo = {'roomName': roomName,'nickName': nickName};
			var changeToken = {'cmd': '2', 'idem': Date.now().toString(), 'payload': loginInfo};
			ws.send(JSON.stringify(changeToken));
		}
		function receivePacketHandle(msg) {
			var packet = JSON.parse(msg)
			if (packet.result == 'err') {
				return;
			}
			switch(packet.cmd){
				case '3':
					// console.log('CMD_R_TOKEN_CHANGE');
					mmtRoomName.textContent = roomName;
					mmtNickName.textContent = nickName;
					mmtUserIcon.src = packet.payload.icon;
					addRoom(roomName);
				break;
				case '5':
					// console.log('CMD_R_PLAYER_LOGOUT');
					ws.close();
					mmtRoomList.innerHTML = '';
				break;
				case '7':
					// console.log('CMD_R_GET_MEMBER_LIST');
					var targetRoomName = packet.idem;
					if (roomName == targetRoomName) {
						mmtRoomUserCount.innerHTML = packet.payload.length;
					}
				break;
				case '9':
					// console.log('CMD_R_PLAYER_EXIT_ROOM');
					var targetRoomName = packet.idem;
					roomMsgList[targetRoomName] = null;
					var li = shadowRoot.querySelector('#room-' + targetRoomName);
					mmtRoomList.removeChild(li);
					if (roomName == targetRoomName) {
						roomName = '';
					}
				break;
				case '11':
					// console.log('CMD_R_PLAYER_ENTER_ROOM');
					var targetRoomName = packet.idem;
					addRoom(targetRoomName);
				break;
				case '13':
					// console.log('CMD_B_PLAYER_ENTER_ROOM');
					var targetRoomName = packet.payload.roomInfo.roomName;
					var targetNickName = packet.payload.from.nickName;
					addSystemMessage(targetRoomName, targetNickName + ' 加入聊天室');
					if (roomName == targetRoomName) {
						userSearch();
					}
				break;
				case '15':
					// console.log('CMD_B_PLAYER_EXIT_ROOM');
					var targetRoomName = packet.payload.roomInfo.roomName;
					var targetNickName = packet.payload.from.nickName;
					addSystemMessage(targetRoomName, targetNickName + ' 離開聊天室');
				break;
				case '17':
					// console.log('CMD_R_GET_ROOM_LIST');
					mmtRoomListPanel.hidden = false;
					addRoomList(packet.payload);
				break;
				case "19":
					// console.log('CMD_R_GET_CHAT_HISTORY');
					var targetRoomName = packet.idem;
					var historyUid = packet.payload.historyUid;
					roomMsgList[targetRoomName].historyUid = historyUid;
					for (var i= 0; i< packet.payload.message.length; i++) {
						var msgObj = packet.payload.message[i];
						addUserMessage(msgObj, true);
					}
				break;
				case '51':
					// console.log('CMD_B_PLAYER_SPEAK');
					var msgObj = packet.payload;
					addUserMessage(msgObj);
				break;
				case '81':
					// console.log('CMD_R_PLAYER_SEND_MSG');
				break;
				default:
			}
			// console.log(packet);
		}
		mmtShowRoomListPanelBtn.onclick = openRoomListPanel;
		mmtJoinBtn.onclick = clickJoinRoom;
		mmtSendBtn.onclick = sendMessage;
		mmtRoomNameInput.onkeydown = function(evt) {
			if (evt.key == 'Enter' && !evt.shiftKey) {
				clickJoinRoom();
				return false;
			}
		};
		mmtMessageInput.onkeydown = function(evt) {
			if (evt.key == 'Enter' && !evt.shiftKey) {
				sendMessage();
				return false;
			}
		};
		connect();
	}
	onLoad();
}();