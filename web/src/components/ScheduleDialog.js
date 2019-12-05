import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import { useTranslation } from 'react-i18next';
import { Route, useHistory, useRouteMatch } from 'react-router';

import ConfigurationDialog from './ConfigurationDialog';
import IntervalList from './IntervalList';
import IntervalDialog from './IntervalDialog';
import { intervalToApi } from '../schedule';
import websocket from '../websocket';

const useStyles = makeStyles(theme => ({
  fab: {
    position: 'absolute',
    bottom: theme.spacing(2),
    right: theme.spacing(2),
  },
}));

export default function ScheduleDialog(props) {
  const { onClose, schedule, outletId } = props;

  const history = useHistory();
  const { path, url } = useRouteMatch();

  const handleEditDialogOpen = (interval) => history.push(`${url}/interval/${interval.id}`);

  const handleCreateDialogOpen = () => history.push(`${url}/interval/new`);

  const handleDialogClose = () => history.push(url);

  const handleToggle = (interval) => {
    interval.enabled = !interval.enabled;

    sendMessage('update', interval);
  }

  const handleDelete = (interval) => sendMessage('delete', interval);

  const handleSave = (interval) => {
    sendMessage(interval.id ? 'update' : 'create', interval);

    handleDialogClose();
  }

  const sendMessage = (action, interval) => {
    const data = {
      action: action,
      id: outletId,
      interval: intervalToApi(interval),
    }

    websocket.sendMessage({ type: 'interval', data });
  }

  const classes = useStyles();
  const { t } = useTranslation();

  return (
    <ConfigurationDialog title={t('schedule')} onClose={onClose}>
      <IntervalList
        intervals={schedule}
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
      <Route path={`${path}/interval/:intervalId`}>
        <IntervalDialog
          onClose={handleDialogClose}
          onDone={handleSave}
          intervals={schedule}
        />
      </Route>
    </ConfigurationDialog>
  );
}

ScheduleDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  schedule: PropTypes.array.isRequired,
  outletId: PropTypes.string.isRequired,
};
