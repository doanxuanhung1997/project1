const domElement = document.querySelector(".chat__app-container");

class App extends React.Component {
    constructor() {
        super();
        this.state = {
            chatUserList: [],
            message: null,
            selectedUserID: null,
            userID: null
        }
        this.webSocketConnection = null;
    }

    componentDidMount() {
        this.setWebSocketConnection();
        this.subscribeToSocketMessage();
    }

    setWebSocketConnection() {
        const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYxN2UxNDJmMDU1ZWUyMDFiNzU1ZDE0NiIsInBob25lX251bWJlciI6IjAzMzUyOTk5MzciLCJlbWFpbCI6Imh1bmdkeEB5b3BtYWlsLmNvbSIsInJvbGUiOjIsImV4cCI6MTY0MDkxMDY0OH0.JWw8rDAA4An8NBUuSS28SuNjOjDGBNRLyFJ-vf1q3T8";
        // const tokenLocal = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjYxOTkyMDM4NzMyNDE2YmMzOTYwMTI0ZCIsInBob25lX251bWJlciI6IjAzMzUyOTk5MzciLCJlbWFpbCI6IiIsInJvbGUiOjEsImV4cCI6MTY0MDMxODMwOX0.X3cmcLqEp71KZiDC7ovGn9AtfaBjlvvt4FY8J_ULtv0";
        if (window["WebSocket"]) {
            const socketConnection = new WebSocket("wss://api.sandexcare.com/ws/" + token);
            // this.webSocketConnection = new WebSocket("ws://" + document.location.host + "/ws/" + tokenLocal);
        }
    }

    subscribeToSocketMessage = () => {
        if (this.webSocketConnection === null) {
            return;
        }

        this.webSocketConnection.onclose = (evt) => {
            this.setState({
                message: 'Your Connection is closed.',
                chatUserList: []
            });
        };

        this.webSocketConnection.onmessage = (event) => {
            try {
                const socketPayload = JSON.parse(event.data);
                console.log(socketPayload);
                switch (socketPayload.eventName) {
                    case 'join':
                    case 'disconnect':
                        if (!socketPayload.eventPayload) {
                            return
                        }

                        const userInitPayload = socketPayload.eventPayload;

                        this.setState({
                            chatUserList: userInitPayload.users,
                            userID: this.state.userID === null ? userInitPayload.userID : this.state.userID
                        });

                        break;
                    default:
                        break;
                }
            } catch (error) {
                console.log(error)
                console.warn('Something went wrong while decoding the Message Payload')
            }
        };
    }

    handleKeyPress = (event) => {
        try {
            if (event.key === 'Enter') {
                if (!this.webSocketConnection) {
                    return false;
                }
                if (!event.target.value) {
                    return false;
                }

                this.webSocketConnection.send(JSON.stringify({
                    EventName: 'destroy_call',
                    EventPayload: {
                        // gửi event thông báo đến chuyên viên
                        receiverId: "617e142f055ee201b755d146",
                        // id của cuộc gọi muốn hủy
                        callId: "111111111111111111111111"
                    },
                }));

                // this.webSocketConnection.send(JSON.stringify({
                //     EventName: 'payment',
                //     EventPayload: {
                //         // gửi event thông báo đến user id
                //         receiverId: "617e142f055ee201b755d146",
                //         // trạng thái thanh toán
                //         status: "success",
                //         // nội dung thanh toán
                //         extraData: "Thanh toán thành công",
                //         // phương thức thanh toán nạp kim cương
                //         method: "momo",
                //         // sô tiền nạp kim cương
                //         total: 1000000
                //     },
                // }));
            
                event.target.value = '';
            }            
        } catch (error) {
            console.log(error)
            console.warn('Something went wrong while decoding the Message Payload')
        }
    }

    setNewUserToChat = (event) => {
        if (event.target && event.target.value) {
            if (event.target.value === "select-user") {
                alert("Select a user to chat");
                return;
            }
            this.setState({
                selectedUserID: event.target.value
            })   
        }
    }
    
    getChatList() {
        if (this.state.chatUserList.length === 0) {
            return(
                <h3>No one has joined yet</h3>
            )
        }
        return (
            <div className="chat__list-container">
                <p>Select a user to chat</p>
                <select onChange={this.setNewUserToChat}>
                    <option value={'select-user'} className="username-list">Select User</option>
                    {
                        this.state.chatUserList.map(user => {
                            if (user.userID !== this.state.userID) {
                                return (
                                    <option value={user.userID} className="username-list">
                                        {user.username}
                                    </option>
                                )
                            }
                        })
                    }
                </select>
            </div>
        );
    }

    getChatContainer() {
        return (
            <div class="chat__message-container">
                <div class="message-container">
                    {this.state.message}
                </div>
                <input type="text" id="message-text" size="64" autofocus placeholder="Type Your message" onKeyPress={this.handleKeyPress}/>
            </div>
        );
    }

    render() {
        return (
            <React.Fragment>
                {this.getChatList()}
                {this.getChatContainer()}
            </React.Fragment>
        );
    }
}

ReactDOM.render(<App />, domElement)