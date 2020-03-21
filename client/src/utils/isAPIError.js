export default err => {
  return (
    err &&
    err.response &&
    err.response.data &&
    Array.isArray(err.response.data.errors) &&
    err.response.data.errors[0]
  );
};
