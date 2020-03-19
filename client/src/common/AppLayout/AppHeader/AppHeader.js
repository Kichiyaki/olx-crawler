import React from 'react';
import classNames from 'classnames';
import { func, bool } from 'prop-types';
import { NAME } from '@config/application';

import { makeStyles } from '@material-ui/core/styles';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Tooltip
} from '@material-ui/core';
import { Menu as MenuIcon, Language as LanguageIcon } from '@material-ui/icons';
import Link from '@common/Link/Link';
import { DRAWER_WIDTH } from '../Drawer/Drawer';

const useStyles = makeStyles(theme => {
  return {
    navItem: {
      marginRight: theme.spacing(2)
    },
    toolbar: {
      paddingRight: 24 // keep right padding when drawer closed
    },
    menuButton: {
      marginLeft: 12,
      marginRight: 36
    },
    menuButtonHidden: {
      display: 'none'
    },
    appBar: ({ noBoxShadow }) => ({
      zIndex: theme.zIndex.drawer + 1,
      transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.leavingScreen
      }),
      boxShadow: noBoxShadow ? 'none' : undefined
    }),
    appBarShift: {
      marginLeft: DRAWER_WIDTH,
      width: `calc(100% - ${DRAWER_WIDTH}px)`,
      transition: theme.transitions.create(['width', 'margin'], {
        easing: theme.transitions.easing.sharp,
        duration: theme.transitions.duration.enteringScreen
      })
    }
  };
});

const AppHeader = ({ isOpen, onDrawerOpen, noBoxShadow }) => {
  const classes = useStyles({ noBoxShadow });
  return (
    <AppBar
      position="absolute"
      className={classNames(classes.appBar, isOpen && classes.appBarShift)}
    >
      <Toolbar disableGutters={!isOpen} className={classes.toolbar}>
        <IconButton
          color="inherit"
          aria-label="Open drawer"
          data-testid="open-drawer"
          onClick={onDrawerOpen}
          className={classNames(
            classes.menuButton,
            isOpen && classes.menuButtonHidden
          )}
        >
          <MenuIcon />
        </IconButton>
        <Typography variant="h4" component="h3" color="inherit" noWrap>
          {NAME}
        </Typography>
        <div style={{ flex: 1 }} />
        <Tooltip title="Strona aplikacji">
          <IconButton color="inherit">
            <Link to="https://dawid-wysokinski.pl">
              <LanguageIcon />
            </Link>
          </IconButton>
        </Tooltip>
      </Toolbar>
    </AppBar>
  );
};

AppHeader.defaultProps = {
  isOpen: false,
  onDrawerOpen: () => {},
  noBoxShadow: false
};

AppHeader.propTypes = {
  onDrawerOpen: func.isRequired,
  isOpen: bool.isRequired,
  noBoxShadow: bool.isRequired
};

export default AppHeader;
