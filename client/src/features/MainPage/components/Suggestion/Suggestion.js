import React from 'react';
import { format } from 'date-fns';
import { DATE_FORMAT } from '@config/application';
import { object } from 'prop-types';

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
    '&:not(:last-child)': {
      marginBottom: theme.spacing(2)
    }
  }
}));

export default function Suggestion({ data }) {
  const classes = useStyles();

  return (
    <Card className={classes.card}>
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
        <Typography component="p">Cena: {data.price}</Typography>
      </CardContent>
      <CardActions disableSpacing className={classes.actions}>
        <Button color="secondary">Usuń</Button>
        <Button>
          <Link color="secondary" underline="none" to={data.url}>
            Przejdź do aukcji
          </Link>
        </Button>
      </CardActions>
    </Card>
  );
}

Suggestion.propTypes = {
  data: object.isRequired
};
