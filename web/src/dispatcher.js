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

  dispatchGroupMessage(groupID, action) {
    this.dispatchMessage('group', { groupID, action });
  }

  dispatchOutletMessage(outletID, action) {
    this.dispatchMessage('outlet', { outletID, action });
  }

  dispatchIntervalMessage(outletID, action, interval) {
    const data = { outletID, action, interval: intervalToApi(interval) };

    this.dispatchMessage('interval', data);
  }

  dispatchMessage(type, data = {}) {
    this.ws.sendMessage({ type, data });
  }
}

const dispatcher = new Dispatcher(websocket);

export default dispatcher;
