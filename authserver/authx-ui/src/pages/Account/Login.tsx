

import { useDispatch, useSelector } from 'react-redux';
import { increment, decrement, set } from '../../components/features/users/UserSlice';
import { RootState } from '../../store/store';

const Login = () => {
  const dispatch = useDispatch();
  const count = useSelector((state: RootState) => state.counter.value);

  return (
    <div>
      <p>Count: {count}</p>
      <button onClick={() => dispatch(increment())}>Increment</button>
      <button onClick={() => dispatch(decrement())}>Decrement</button>
      <button onClick={() => dispatch(set(0))}>Reset</button>
    </div>
  );
  }
  
  export default Login;