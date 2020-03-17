import React from 'react';
import { string as propString, node } from 'prop-types';
import { Link } from 'react-router-dom';
import { Link as MaterialUILink } from '@material-ui/core';

const MyLink = ({ children, to, ...rest }) => {
  return (
    <MaterialUILink
      component={to.includes('http') ? undefined : Link}
      to={to}
      underline="none"
      color="inherit"
      {...rest}
    >
      {children}
    </MaterialUILink>
  );
};

MyLink.propTypes = {
  to: propString.isRequired,
  children: node
};

export default MyLink;
