# Testing Guide

This document describes the testing setup for the VGC Corner frontend application.

## Testing Stack

- **Jest**: Test runner and assertion library
- **React Testing Library**: Component testing utilities
- **@testing-library/jest-dom**: Custom Jest matchers for DOM assertions

## Running Tests

```bash
# Run all tests once
npm test

# Run tests in watch mode (reruns on file changes)
npm run test:watch

# Run tests with coverage report
npm run test:coverage
```

## Test Structure

Tests are organized using the following conventions:

- Test files use the `.test.tsx` or `.test.ts` extension
- Tests are located in `__tests__` directories adjacent to the code they test
- Example structure:
  ```
  src/
    components/
      battles/
        BattleHeader.tsx
        __tests__/
          BattleHeader.test.tsx
  ```

## Writing Tests

### Component Tests

Component tests should verify:
- Rendering with different props
- User interactions
- Conditional rendering logic
- Accessibility

Example:

```typescript
import { render, screen } from '@testing-library/react';
import MyComponent from '../MyComponent';

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent prop="value" />);
    expect(screen.getByText('Expected Text')).toBeInTheDocument();
  });
});
```

### Mocking Next.js Components

Next.js components like `Image` and `Link` need to be mocked in tests:

```typescript
jest.mock('next/image', () => ({
  __esModule: true,
  default: (props: any) => <img {...props} />,
}));
```

## Test Coverage

Current test coverage includes:

1. **Home Page** (`src/app/__tests__/page.test.tsx`)
   - Renders main content
   - Contains proper links
   - Shows logos

2. **BattleHeader Component** (`src/components/battles/__tests__/BattleHeader.test.tsx`)
   - Displays player information
   - Shows winner badge
   - Formats battle metadata
   - Shows statistics

3. **TeamComparison Component** (`src/components/battles/__tests__/TeamComparison.test.tsx`)
   - Renders team members
   - Shows HP bars
   - Displays status conditions
   - Shows abilities and items
   - Indicates fainted Pokemon

## Best Practices

1. **Use semantic queries**: Prefer `getByRole`, `getByLabelText`, etc. over `getByTestId`
2. **Test user behavior**: Focus on what users see and do, not implementation details
3. **Keep tests simple**: One concept per test
4. **Use meaningful test names**: Describe what the test verifies
5. **Clean up**: Tests should be independent and not affect each other

## Configuration

Testing configuration is located in:
- `jest.config.js` - Main Jest configuration
- `jest.setup.ts` - Test environment setup (imports jest-dom matchers)

## Common Issues

### "Cannot find module" errors
Make sure the module path alias in `jest.config.js` matches your `tsconfig.json`:
```javascript
moduleNameMapper: {
  '^@/(.*)$': '<rootDir>/src/$1',
}
```

### Image/Link component errors
Mock Next.js components that don't work well in the test environment (see "Mocking Next.js Components" above).

## Adding New Tests

When adding new components:

1. Create a `__tests__` directory next to your component
2. Create a test file matching the component name: `ComponentName.test.tsx`
3. Write tests covering the main functionality
4. Run tests to verify they pass
5. Check coverage with `npm run test:coverage`
