import WebSocketAsPromised from 'websocket-as-promised';

export default class WebSocket {
  defaultListeners = {
    onOpen: (event) => {
      console.log(`Opened websocket connection`, event);
    },
    onClose: (event) => {
      console.log(`Closed websocket connection`, event);
    },
    onError: (event) => {
      console.error(`Websocket error`, event);
    },
  };

  constructor(url) {
    this.ws = new WebSocketAsPromised(url, {
      packMessage: data => JSON.stringify(data),
      unpackMessage: message => JSON.parse(message),
    });
  }

  sendMessage = (data) => this.ws.open().then(() => this.ws.sendPacked(data));

  onMessage = (cb) => this.ws.onUnpackedMessage.addListener(cb);

  onOpen = (cb) => this.ws.onOpen.addListener(cb);

  onClose = (cb) => this.ws.onClose.addListener(cb);

  onError = (cb) => this.ws.onError.addListener(cb);

  attachDefaultListeners = () => {
    this.onOpen(this.defaultListeners.onOpen);
    this.onClose(this.defaultListeners.onClose);
    this.onError(this.defaultListeners.onError);
  }
}
