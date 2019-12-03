import React, { useReducer } from 'react';
import PropTypes from 'prop-types';
import { useTranslation } from 'react-i18next';
import CheckIcon from '@material-ui/icons/Check';

import ConfigurationDialog from './ConfigurationDialog';
import IntervalOptionsList from './IntervalOptionsList';
import IntervalTimePicker from './IntervalTimePicker';
import WeekdaysDialog from './WeekdaysDialog';

const reduceState = (state, changes) => ({ ...state, ...changes });

const initState = (initialState) => {
  const dialogState = {
    fromOpen: false,
    toOpen: false,
    weekdaysOpen: false
  }

  return { ...initialState, ...dialogState };
}

export default function IntervalDialog(props) {
  const { open, onClose, onDone, interval } = props;

  const [state, setState] = useReducer(reduceState, interval, initState);

  const handleOpen = (name, open) => () => {
    setState({ [name + 'Open']: open });
  }

  const handleChange = (name) => (value) => {
    setState({ [name]: value, [name + 'Open']: false });
  }

  const handleDone = () => {
    const newInterval = {
      ...interval,
      weekdays: state.weekdays,
      from: state.from,
      to: state.to
    };

    onDone(newInterval);
  };

  const isComplete = () => state.from && state.to && state.weekdays.length > 0;

  const { t } = useTranslation();

  return (
    <ConfigurationDialog
      title={state.id ? t('edit-interval') : t('add-interval')}
      open={open}
      onClose={onClose}
      onDone={handleDone}
      doneButtonDisabled={!isComplete()}
      doneButtonText={<CheckIcon />}
    >
      <IntervalOptionsList
        weekdays={state.weekdays}
        fromDayTime={state.from}
        toDayTime={state.to}
        onWeekdaysClick={handleOpen('weekdays', true)}
        onFromDayTimeClick={handleOpen('from', true)}
        onToDayTimeClick={handleOpen('to', true)}
      />
      <IntervalTimePicker
        open={state.fromOpen}
        value={state.from}
        onChange={handleChange('from')}
        onClose={handleOpen('from', false)}
      />
      <IntervalTimePicker
        open={state.toOpen}
        value={state.to}
        onChange={handleChange('to')}
        onClose={handleOpen('to', false)}
      />
      <WeekdaysDialog
        open={state.weekdaysOpen}
        selected={state.weekdays}
        key={interval.id}
        onDone={handleChange('weekdays')}
        onClose={handleOpen('weekdays', false)}
      />
    </ConfigurationDialog>
  );
}

IntervalDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  interval: PropTypes.object.isRequired,
  onClose: PropTypes.func.isRequired,
  onDone: PropTypes.func.isRequired,
};
