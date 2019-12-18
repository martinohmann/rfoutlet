import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import SettingsIcon from '@material-ui/icons/Settings';
import { NoItemsListItem } from '../List';
import { useTranslation } from 'react-i18next';
import { useHistory } from 'react-router';
import { formatSchedule } from '../../format';
import dispatcher from '../../dispatcher';

const useStyles = makeStyles(theme => ({
  container: {
    padding: 0,
  },
}));

export default function OutletList({ outlets }) {
  const classes = useStyles();
  const { t } = useTranslation();

  return (
    <List className={classes.container}>
      {outlets.map(outlet =>
        <OutletListItem
          key={outlet.id}
          {...outlet}
          schedule={outlet.schedule}
        />
      )}
      {outlets.length === 0 && (
        <NoItemsListItem
          primary={t('No outlets configured for this group.')}
          secondary={t('Check your rfoutlet config.')}
        />
      )}
    </List>
  );
}

OutletList.propTypes = {
  outlets: PropTypes.array.isRequired,
};


const OutletListItem = ({ id, name, state, schedule }) => {
  const history = useHistory();
  const { t } = useTranslation();

  const handleToggle = () => dispatcher.dispatchOutletMessage(id, 'toggle');

  const hasEnabledIntervals = () => schedule.some(interval => interval.enabled);

  return (
    <ListItem>
      <ListItemText
        primary={name}
        secondary={formatSchedule(schedule, t)}
        onClick={() => history.push(`/schedule/${id}`)}
      />
      <ListItemSecondaryAction>
        <Switch
          color="primary"
          onChange={handleToggle}
          checked={state === 1}
          disabled={hasEnabledIntervals()}
        />
        <IconButton onClick={() => history.push(`/schedule/${id}`)}>
          <SettingsIcon />
        </IconButton>
      </ListItemSecondaryAction>
    </ListItem>
  );
};

OutletListItem.propTypes = {
  id: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  state: PropTypes.number.isRequired,
  schedule: PropTypes.array.isRequired,
};
