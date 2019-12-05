import React, { useReducer } from 'react';
import PropTypes from 'prop-types';
import { useTranslation } from 'react-i18next';
import CheckIcon from '@material-ui/icons/Check';
import { Route, Switch, useHistory, useRouteMatch } from 'react-router';

import ConfigurationDialog from './ConfigurationDialog';
import IntervalOptionsList from './IntervalOptionsList';
import IntervalTimePicker from './IntervalTimePicker';
import WeekdaysDialog from './WeekdaysDialog';

const emptyInterval = {
  id: null,
  enabled: false,
  from: null,
  to: null,
  weekdays: [],
}

const findInterval = (intervals, id) => {
  const interval = intervals.find(interval => interval.id === id);
  if (interval === undefined) {
    return emptyInterval;
  }

  return interval;
}

const reduceState = (state, changes) => ({ ...state, ...changes });

export default function IntervalDialog(props) {
  const { onClose, onDone, intervals } = props;

  const history = useHistory();
  const { path, url, params } = useRouteMatch();

  const interval = findInterval(intervals, params.intervalId);

  const [state, setState] = useReducer(reduceState, interval);

  const handleOpen = (name, open) => () => history.push(open ? `${url}/${name}` : url);

  const handleChange = (name) => (value) => setState({ [name]: value });

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
      <Switch>
        <Route path={`${path}/from`}>
          <IntervalTimePicker
            value={state.from}
            onChange={handleChange('from')}
            onClose={handleOpen('from', false)}
          />
        </Route>
        <Route path={`${path}/to`}>
          <IntervalTimePicker
            value={state.to}
            onChange={handleChange('to')}
            onClose={handleOpen('to', false)}
          />
        </Route>
        <Route path={`${path}/weekdays`}>
          <WeekdaysDialog
            selected={state.weekdays}
            key={interval.id}
            onChange={handleChange('weekdays')}
            onClose={handleOpen('weekdays', false)}
          />
        </Route>
      </Switch>
    </ConfigurationDialog>
  );
}

IntervalDialog.propTypes = {
  onClose: PropTypes.func.isRequired,
  onDone: PropTypes.func.isRequired,
  intervals: PropTypes.array.isRequired,
};
