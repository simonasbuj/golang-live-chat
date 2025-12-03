function chatApp() {
    return {
        inputUsername: "",
        inputChatRoomCode: "",
        inputMessage: "",

        username: "",
        chatRoomCode: "",

        socket: null,
        messages: [],

        async init() {
            console.log("chatApp.js is LOADED")
        },
        
        join() {
            if (!this.inputChatRoomCode.trim() || !this.inputUsername.trim()) return

            this.username = this.inputUsername
            this.chatRoomCode = this.inputChatRoomCode

            const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
            const host = window.location.host

            this.socket = new WebSocket(`${protocol}://${host}/chat?room=${this.chatRoomCode}&token=${this.username}`)

            this.socket.addEventListener("open", () => {
                console.log("connected to chat: ", this.chatRoomCode)
            })

            this.socket.addEventListener("message", (event) => this.handleReceiveMsg(event));

            this.socket.addEventListener("close", () => {
                console.log("Disconnected. Reconnecting in 2s…");
                setTimeout(() => this.connect(), 2000);
            });

            this.inputUsername = ""
            this.inputChatRoomCode = ""
        },

        handleReceiveMsg(event) {
            try {
                const parsedMsg =JSON.parse(event.data)
                this.messages.push(parsedMsg)
                this.scrollToTextBoxBottom()
            } catch (e) {
                console.error("couldn't parse message data: ", event.data, " error: ", e)
            }

            if (this.messages.length > 200) {
                this.messages.shift()
            }
        },

        send() {
            if (!this.inputMessage.trim()) return

            const msg = {
                user: this.username,
                content: this.inputMessage,
            }

            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
                this.socket.send(JSON.stringify(msg));
            }

            this.inputMessage = ""
        },

        scrollToTextBoxBottom() {
            this.$nextTick(() => {
                const el = this.$refs.chatbox;
                el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' });
            });
        },
    }
}
