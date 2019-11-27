import WebSocketAsPromised from 'websocket-as-promised';
import config from './config';

export class Websocket {
  defaultListeners = {
    onOpen: (event) => console.log(`Opened websocket connection`, event),
    onClose: (event) => console.log(`Closed websocket connection`, event),
    onError: (event) => console.error(`Websocket error`, event),
    onMessage: (msg) => console.log("[ws recv]", msg),
  };

  constructor(url) {
    this.ws = new WebSocketAsPromised(url, {
      packMessage: (data) => JSON.stringify(data),
      unpackMessage: (message) => JSON.parse(message),
    });
  }

  sendMessage = (data) => this.ws.open()
    .then(() => {
      console.log("[ws send]", data);

      this.ws.sendPacked(data);
    })
    .catch(err => console.error(err));

  onMessage = (cb) => this.ws.onUnpackedMessage.addListener(cb);

  onOpen = (cb) => this.ws.onOpen.addListener(cb);

  onClose = (cb) => this.ws.onClose.addListener(cb);

  onError = (cb) => this.ws.onError.addListener(cb);

  attachDefaultListeners = () => {
    this.onOpen(this.defaultListeners.onOpen);
    this.onClose(this.defaultListeners.onClose);
    this.onError(this.defaultListeners.onError);
    this.onMessage(this.defaultListeners.onMessage);
  }
}

const websocket = new Websocket(config.ws.url);

websocket.attachDefaultListeners();

export default websocket;
