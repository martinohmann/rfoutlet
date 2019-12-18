import React from 'react';
import PropTypes from 'prop-types';
import { List, ListItem, NoItemsListItem } from '../List';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Switch from '@material-ui/core/Switch';
import { useTranslation } from 'react-i18next';
import IntervalActionsMenu from './IntervalActionsMenu';
import { formatDayTimeInterval, formatWeekdays } from '../../format';

export default function IntervalList({ intervals, onDelete, onEdit, onToggle }) {
  const { t } = useTranslation();

  return (
    <List>
      {intervals.map(interval => (
        <IntervalListItem
          key={interval.id}
          interval={interval}
          onToggle={() => onToggle(interval)}
          onEdit={() => onEdit(interval)}
          onDelete={() => onDelete(interval)}
        />
      ))}
      {intervals.length === 0 && (
        <NoItemsListItem
          primary={t('no-intervals-primary')}
          secondary={t('no-intervals-secondary')}
        />
      )}
    </List>
  );
}

IntervalList.propTypes = {
  intervals: PropTypes.array.isRequired,
  onDelete: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onToggle: PropTypes.func.isRequired,
};

const IntervalListItem = ({ interval, onDelete, onEdit, onToggle }) => {
  const { t } = useTranslation();

  return (
    <ListItem onClick={onEdit}>
      <ListItemText
        primary={formatDayTimeInterval(interval)}
        secondary={formatWeekdays(interval.weekdays, t)}
      />
      <ListItemSecondaryAction>
        <Switch
          color="primary"
          checked={interval.enabled}
          onChange={onToggle}
        />
        <IntervalActionsMenu onEdit={onEdit} onDelete={onDelete} />
      </ListItemSecondaryAction>
    </ListItem>
  );
};

IntervalListItem.propTypes = {
  interval: PropTypes.object,
  onDelete: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onToggle: PropTypes.func.isRequired,
};
