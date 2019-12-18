import websocket from './websocket';
import { convertToApp, intervalToApi } from './convert';

class Dispatcher {
  constructor(websocket) {
    this.ws = websocket;
  }

  addMessageListener(listener) {
    this.ws.onMessage(msg => listener(convertToApp(msg)));
  }

  dispatchStatusMessage() {
    this.dispatchMessage('status');
  }

  dispatchGroupMessage(id, action) {
    this.dispatchMessage('group', { id, action });
  }

  dispatchOutletMessage(id, action) {
    this.dispatchMessage('outlet', { id, action });
  }

  dispatchIntervalMessage(id, action, interval) {
    const data = { id, action, interval: intervalToApi(interval) };

    this.dispatchMessage('interval', data);
  }

  dispatchMessage(type, data = {}) {
    this.ws.sendMessage({ type, data });
  }
}

const dispatcher = new Dispatcher(websocket);

export default dispatcher;
