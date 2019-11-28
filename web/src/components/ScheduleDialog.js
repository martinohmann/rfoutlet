import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';

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

const emptyInterval = {
  id: null,
  enabled: false,
  from: null,
  to: null,
  weekdays: [],
}

export default function ScheduleDialog(props) {
  const { open, onClose, schedule, outletId } = props;

  const [dialogOpen, setDialogOpen] = useState(false);
  const [currentInterval, setCurrentInterval] = useState(emptyInterval);

  const handleDialogOpen = (interval) => {
    setDialogOpen(true);
    setCurrentInterval(interval);
  }

  const handleDialogClose = () => {
    setDialogOpen(false);
    setCurrentInterval(emptyInterval);
  }

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

  return (
    <ConfigurationDialog title="Schedule" open={open} onClose={onClose}>
      <IntervalList
        intervals={schedule}
        onToggle={handleToggle}
        onEdit={handleDialogOpen}
        onDelete={handleDelete}
      />
      <Fab
        color="secondary"
        className={classes.fab}
        onClick={() => handleDialogOpen(emptyInterval)}
      >
        <AddIcon />
      </Fab>
      <IntervalDialog
        open={dialogOpen}
        onClose={handleDialogClose}
        onDone={handleSave}
        key={currentInterval.id}
        interval={currentInterval}
      />
    </ConfigurationDialog>
  );
}

ScheduleDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  schedule: PropTypes.array.isRequired,
  outletId: PropTypes.string.isRequired,
};
