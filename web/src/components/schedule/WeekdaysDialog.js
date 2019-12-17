import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { List, ListItem } from '../List';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import Checkbox from '@material-ui/core/Checkbox';
import { useTranslation } from 'react-i18next';
import CheckIcon from '@material-ui/icons/Check';
import Dialog from '../Dialog';
import { weekdaysLong } from '../../format';

export default function WeekdaysDialog(props) {
  const { onClose, onChange } = props;

  const [selected, setSelected] = useState([...props.selected]);

  const handleWeekdayToggle = (key) => () => {
    const selectedIndex = selected.indexOf(key);
    let newSelected = [];

    if (selectedIndex === -1) {
      newSelected = newSelected.concat(selected, key);
    } else if (selectedIndex === 0) {
      newSelected = newSelected.concat(selected.slice(1));
    } else if (selectedIndex === selected.length - 1) {
      newSelected = newSelected.concat(selected.slice(0, -1));
    } else if (selectedIndex > 0) {
      newSelected = newSelected.concat(
        selected.slice(0, selectedIndex),
        selected.slice(selectedIndex + 1),
      );
    }

    setSelected(newSelected);
  }

  const isSelected = (key) => selected.indexOf(key) !== -1;

  const handleDone = () => {
    onChange(selected);
    onClose();
  }

  const { t } = useTranslation();

  return (
    <Dialog
      title={t('select-weekdays')}
      onClose={onClose}
      onDone={handleDone}
      doneButtonDisabled={selected.length === 0}
      doneButtonText={<CheckIcon />}
    >
      <List>
        {weekdaysLong.map((weekday, key) => (
          <WeekdayListItem
            key={key}
            weekday={weekday}
            selected={isSelected(key)}
            onToggle={handleWeekdayToggle(key)}
          />
        ))}
      </List>
    </Dialog>
  );
}

WeekdaysDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  onChange: PropTypes.func.isRequired,
};


const WeekdayListItem = ({ onToggle, weekday, selected }) => {
  const { t } = useTranslation();

  return (
    <ListItem onClick={onToggle}>
      <ListItemIcon>
        <Checkbox
          color="primary"
          onChange={onToggle}
          checked={selected}
        />
      </ListItemIcon>
      <ListItemText primary={t(weekday)} />
    </ListItem>
  );
};

WeekdayListItem.propTypes = {
  onToggle: PropTypes.func.isRequired,
  selected: PropTypes.bool.isRequired,
  weekday: PropTypes.string.isRequired,
};
