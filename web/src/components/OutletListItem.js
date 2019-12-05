import React from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import SettingsIcon from '@material-ui/icons/Settings';
import { useTranslation } from 'react-i18next';
import { Route, useHistory } from 'react-router';

import ScheduleDialog from './ScheduleDialog';
import { formatSchedule } from '../schedule';
import websocket from '../websocket';

export default function OutletListItem(props) {
  const { id, name, state, schedule } = props;

  const history = useHistory();

  const handleDialogOpen = (open) => () => history.push(open ? `/schedule/${id}` : '/');

  const handleToggle = () => {
    websocket.sendMessage({ type: 'outlet', data: { id, action: 'toggle' } });
  }

  const hasEnabledIntervals = () => schedule.some(interval => interval.enabled);

  const { t } = useTranslation();

  return (
    <ListItem>
      <ListItemText primary={name} secondary={formatSchedule(schedule, t)} onClick={handleDialogOpen(true)} />
      <Route path={`/schedule/${id}`}>
        <ScheduleDialog
          outletId={id}
          schedule={schedule}
          onClose={handleDialogOpen(false)}
        />
      </Route>
      <ListItemSecondaryAction>
        <Switch
          color="primary"
          onChange={handleToggle}
          checked={state === 1}
          disabled={hasEnabledIntervals()}
        />
        <IconButton onClick={handleDialogOpen(true)}>
          <SettingsIcon />
        </IconButton>
      </ListItemSecondaryAction>
    </ListItem>
  );
}
