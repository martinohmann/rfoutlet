export default {
  project: {
    name: "rfoutlet",
    url: 'https://github.com/martinohmann/rfoutlet',
  },
  api: {
    baseUri: (() => {
      if (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') {
        const l = window.location;

        return `${l.protocol}//${l.hostname}:3333/api`
      }

      return '/api';
    })(),
  }
}
