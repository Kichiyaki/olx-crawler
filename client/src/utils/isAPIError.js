export default err => {
  return (
    err &&
    err.response &&
    err.response.data &&
    err.response.data.data &&
    Array.isArray(err.response.data.data.errors) &&
    err.response.data.data.errors[0]
  );
};
