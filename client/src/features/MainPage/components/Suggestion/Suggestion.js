import React from 'react';
import { object, func, bool } from 'prop-types';
import { format } from 'date-fns';
import classnames from 'classnames';
import { DATE_FORMAT } from '@config/application';

import { darken } from '@material-ui/core/styles/colorManipulator';
import { makeStyles } from '@material-ui/core/styles';
import {
  Card,
  CardHeader,
  CardMedia,
  CardContent,
  CardActions,
  Typography,
  Button,
  Hidden,
  Box
} from '@material-ui/core';
import Link from '@common/Link/Link';
import { red } from '@material-ui/core/colors';

const useStyles = makeStyles(theme => ({
  media: {
    height: 0,
    paddingTop: '56.25%' // 16:9
  },
  avatar: {
    backgroundColor: red[500]
  },
  actions: {
    justifyContent: 'flex-end'
  },
  card: {
    transition: 'all .2s'
  },
  selected: {
    backgroundColor: darken(theme.palette.background.paper, 0.1)
  }
}));

export default function Suggestion({ data, onSelect, selected, t }) {
  const classes = useStyles({ selected });

  return (
    <Card
      className={classnames(classes.card, { [classes.selected]: selected })}
    >
      <CardHeader
        action={
          data.observation ? (
            <Hidden xsDown>
              <Box display="flex" mt={1} mr={1}>
                <Link color="secondary" to={data.observation.url}>
                  {data.observation.name}
                </Link>
              </Box>
            </Hidden>
          ) : (
            undefined
          )
        }
        title={data.title}
        subheader={format(new Date(data.created_at), DATE_FORMAT)}
      />
      {data.image && (
        <CardMedia
          className={classes.media}
          image={data.image}
          title={data.title}
        />
      )}
      <CardContent>
        <Typography component="p">
          {t('price', { value: data.price })}
        </Typography>
      </CardContent>
      <CardActions disableSpacing className={classes.actions}>
        <Button onClick={onSelect} color="secondary">
          {t('choose')}
        </Button>
        <Button>
          <Link color="secondary" underline="none" to={data.url}>
            {t('goToAuction')}
          </Link>
        </Button>
      </CardActions>
    </Card>
  );
}

Suggestion.propTypes = {
  data: object.isRequired,
  t: func.isRequired,
  onSelect: func.isRequired,
  selected: bool.isRequired
};
