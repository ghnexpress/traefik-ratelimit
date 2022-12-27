package sliding_window_counter

import "context"

func (r *repository) RemoveExpiredWindowSlice(ctx context.Context, ip string, currSlice, windowTime int) (err error) {
	allReqCount, err := r.GetRequestCountByIP(ctx, ip)
	if err != nil {
		return err
	}
	partToBeEvicted := make([]int, 0)
	for part, _ := range allReqCount {
		if part < currSlice-windowTime {
			partToBeEvicted = append(partToBeEvicted, part)
		}
	}

	if len(partToBeEvicted) < 0 {
		return nil
	}

	for _, part := range partToBeEvicted {
		delete(allReqCount, part)
	}
	r.memory.CompareAndSwap()
	return nil
}
