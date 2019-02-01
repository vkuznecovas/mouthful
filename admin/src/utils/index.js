export const resolveURL = () => {
  // NOTE: https://medium.com/@trekinbami/using-environment-variables-in-react-6b0a99d83cf5

  if (process.env.NODE_ENV === 'production') {
    return `http://${window.location.origin}`;
  } else {
    return `http://${process.env.PREACT_APP_URL}:${process.env.PREACT_APP_PORT}`
  }
}

export const axiosConfig = {
  baseURL: resolveURL(),
};