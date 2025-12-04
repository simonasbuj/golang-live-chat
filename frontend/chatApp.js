function chatApp() {
    return {
        inputUsername: "",
        inputChatRoomCode: "",
        inputMessage: "",

        username: "",
        chatRoomCode: null,

        socket: null,
        messages: [],

        isLoading: false,

        async init() {
            console.log("chatApp.js is LOADED")
        },
        
        join() {
            if (!this.inputChatRoomCode.trim() || !this.inputUsername.trim()) return

            this.isLoading = true
            this.messages = []
            this.username = this.inputUsername
            this.chatRoomCode = this.inputChatRoomCode

            const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
            const host = window.location.host

            this.socket = new WebSocket(`${protocol}://${host}/chat?room=${this.chatRoomCode}&token=${this.username}`)

            this.socket.addEventListener("open", () => {
                console.log("connected to chat: ", this.chatRoomCode)
            })

            this.socket.addEventListener("message", (event) => this.handleReceiveMsg(event));

            this.socket.addEventListener("close", (event) => {
                if (event.code === 1000) {
                    console.log("Normal closure — not reconnecting.");
                    return;
                }

                console.log("Unexpected disconnect. Reconnecting in 2s…");
                setTimeout(() => this.join(), 2000);
            });

            this.isLoading = false
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
                content: this.inputMessage,
            }

            if (this.socket && this.socket.readyState === WebSocket.OPEN) {
                this.socket.send(JSON.stringify(msg));
            }

            this.inputMessage = ""
        },

        leave() {
            this.socket.close(1000, "User left the room")
            this.socket = null

            this.chatRoomCode = null
            this.messages = []
        },

        scrollToTextBoxBottom() {
            this.$nextTick(() => {
                const el = this.$refs.chatbox;
                el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' });
            });
        },
    }
}
