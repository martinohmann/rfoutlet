import React from 'react';
import PropTypes from 'prop-types';
import { ListItem } from './List';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import Checkbox from '@material-ui/core/Checkbox';
import { useTranslation } from 'react-i18next';

export default function WeekdayListItem(props) {
  const { onToggle, weekday, selected } = props;
  const { t } = useTranslation();

  return (
    <ListItem onClick={onToggle}>
      <ListItemText primary={t(weekday)} />
      <ListItemSecondaryAction>
        <Checkbox
          color="primary"
          onChange={onToggle}
          checked={selected}
        />
      </ListItemSecondaryAction>
    </ListItem>
  );
}

WeekdayListItem.propTypes = {
  onToggle: PropTypes.func.isRequired,
  selected: PropTypes.bool.isRequired,
  weekday: PropTypes.string.isRequired,
};
