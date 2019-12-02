import React from 'react';
import PropTypes from 'prop-types';
import { List, NoItemsListItem } from './List';
import IntervalListItem from './IntervalListItem';
import { useTranslation } from 'react-i18next';

export default function IntervalList(props) {
  const { intervals, onDelete, onEdit, onToggle } = props;
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
      {intervals.length === 0 ? (
        <NoItemsListItem
          primary={t('no-intervals-primary')}
          secondary={t('no-intervals-secondary')}
        />
      ) : ''}
    </List>
  );
}

IntervalList.propTypes = {
  intervals: PropTypes.array.isRequired,
  onDelete: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onToggle: PropTypes.func.isRequired,
};
