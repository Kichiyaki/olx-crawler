import React from 'react';
import { string as propString, node } from 'prop-types';

import { makeStyles } from '@material-ui/core/styles';
import { ListItem, ListItemIcon, ListItemText } from '@material-ui/core';
import Link from '@common/Link/Link';

const useStyles = makeStyles(() => ({
  link: {
    textDecoration: 'none'
  }
}));

const MyListItem = ({ to, text, children }) => {
  const classes = useStyles();
  return (
    <Link to={to} underline="none" color="inherit">
      <ListItem button className={classes.link}>
        <ListItemIcon>{children}</ListItemIcon>
        <ListItemText primary={text} />
      </ListItem>
    </Link>
  );
};

MyListItem.propTypes = {
  to: propString.isRequired,
  text: propString.isRequired,
  children: node
};

export default MyListItem;
