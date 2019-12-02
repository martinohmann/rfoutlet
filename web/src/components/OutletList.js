import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import List from '@material-ui/core/List';
import { NoItemsListItem } from './List';
import { useTranslation } from 'react-i18next';

import OutletListItem from './OutletListItem';
import { scheduleToApp } from '../schedule';

const useStyles = makeStyles(theme => ({
  container: {
    padding: 0,
  },
}));

export default function OutletList(props) {
  const classes = useStyles();

  const { outlets } = props;
  const { t } = useTranslation();

  return (
    <List className={classes.container}>
      {outlets.map(outlet =>
        <OutletListItem
          key={outlet.id}
          {...outlet}
          schedule={scheduleToApp(outlet.schedule)}
        />
      )}
      {outlets.length === 0 ? (
        <NoItemsListItem
          primary={t('No outlets configured for this group.')}
          secondary={t('Check your rfoutlet config.')}
        />
      ) : ''}
    </List>
  );
}

OutletList.propTypes = {
  outlets: PropTypes.array.isRequired,
};
