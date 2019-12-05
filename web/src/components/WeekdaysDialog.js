import React, { useState } from 'react';
import PropTypes from 'prop-types';
import { List } from './List';
import { useTranslation } from 'react-i18next';
import CheckIcon from '@material-ui/icons/Check';

import ConfigurationDialog from './ConfigurationDialog';
import WeekdayListItem from './WeekdayListItem';
import { weekdaysLong } from '../schedule';

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
    <ConfigurationDialog
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
    </ConfigurationDialog>
  );
}

WeekdaysDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  onChange: PropTypes.func.isRequired,
};
