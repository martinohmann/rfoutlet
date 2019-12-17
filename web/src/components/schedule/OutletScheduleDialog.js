import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import { useTranslation } from 'react-i18next';
import { Redirect, useHistory, useRouteMatch } from 'react-router';
import IntervalList from './IntervalList';
import Dialog from '../Dialog';
import { useCurrentOutlet } from '../../hooks';
import dispatcher from '../../dispatcher';

const useStyles = makeStyles(theme => ({
  fab: {
    position: 'absolute',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
  },
}));

export default function OutletScheduleDialog({ onClose }) {
  const history = useHistory();
  const { url } = useRouteMatch();

  const outlet = useCurrentOutlet();

  const handleEditDialogOpen = (interval) => history.push(`${url}/interval/${interval.id}`);

  const handleCreateDialogOpen = () => history.push(`${url}/interval/new`);

  const handleToggle = (interval) => {
    interval.enabled = !interval.enabled;

    dispatchMessage('update', interval);
  }

  const handleDelete = (interval) => dispatchMessage('delete', interval);

  const dispatchMessage = (action, interval) => {
    dispatcher.dispatchIntervalMessage(outlet.id, action, interval);
  }

  const classes = useStyles();
  const { t } = useTranslation();

  if (!outlet) {
    return <Redirect to="/" />;
  }

  return (
    <Dialog title={t('schedule')} onClose={onClose}>
      <IntervalList
        intervals={outlet.schedule}
        onToggle={handleToggle}
        onEdit={handleEditDialogOpen}
        onDelete={handleDelete}
      />
      <Fab
        color="secondary"
        className={classes.fab}
        onClick={handleCreateDialogOpen}
      >
        <AddIcon />
      </Fab>
    </Dialog>
  );
}

OutletScheduleDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
};
