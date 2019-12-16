import { useContext } from 'react';
import Context from './Context';
import { useParams } from 'react-router';

export function useOutlets() {
  return useContext(Context).outlets;
}

export function useOutlet(outletId) {
  const outlets = useOutlets();

  return outlets.find(outlet => outlet.id === outletId);
}

export function useCurrentOutlet() {
  const { outletId } = useParams();

  return useOutlet(outletId);
}

export function useIntervals() {
  return useContext(Context).intervals;
}

export function useInterval(intervalId) {
  const intervals = useIntervals();

  return intervals.find(interval => interval.id === intervalId);
}

export function useCurrentInterval() {
  const { intervalId } = useParams();

  return useInterval(intervalId);
}


