import React, { useState } from 'react';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import IconButton from '@material-ui/core/IconButton';
import Switch from '@material-ui/core/Switch';
import EditIcon from '@material-ui/icons/Edit';
import { useTranslation } from 'react-i18next';

import ScheduleDialog from './ScheduleDialog';
import { formatSchedule } from '../schedule';
import websocket from '../websocket';

export default function OutletListItem(props) {
  const { id, name, state, schedule } = props;

  const [scheduleDialogOpen, setScheduleDialogOpen] = useState(false);

  const handleDialogOpen = (open) => () => setScheduleDialogOpen(open);

  const handleToggle = () => {
    websocket.sendMessage({ type: 'outlet', data: { id, action: 'toggle' } });
  }

  const hasEnabledIntervals = () => schedule.some(interval => interval.enabled);

  const { t } = useTranslation();

  return (
    <ListItem>
      <ListItemText primary={name} secondary={formatSchedule(schedule, t)} />
      <ScheduleDialog
        outletId={id}
        schedule={schedule}
        open={scheduleDialogOpen}
        onClose={handleDialogOpen(false)}
      />
      <ListItemSecondaryAction>
        <Switch
          color="primary"
          onChange={handleToggle}
          checked={state === 1}
          disabled={hasEnabledIntervals()}
        />
        <IconButton onClick={handleDialogOpen(true)}>
          <EditIcon />
        </IconButton>
      </ListItemSecondaryAction>
    </ListItem>
  );
}
