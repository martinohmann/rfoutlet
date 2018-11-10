import { api } from './config';

export function makeApiRequest(requestUri, data, success) {
  fetch(api.baseUri + requestUri, {
    method: "POST",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(data)
  }).then(response => {
    return response.json();
  }).then(result => {
      return success(result);
  }).catch(err => {
    console.log(err);
  });
}

export function isOutletEnabled(outlet) {
  if (undefined === outlet || undefined === outlet.state) {
    return false;
  }

  return 1 === outlet.state;
}
