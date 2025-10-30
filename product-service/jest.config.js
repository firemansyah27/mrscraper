/** @type {import('ts-jest').JestConfigWithTsJest} */
module.exports = {
    preset: 'ts-jest',
    testEnvironment: 'node',
    moduleFileExtensions: ['ts', 'js', 'json'],
    rootDir: '.',
    testRegex: '.*\\.spec\\.ts$',
    globals: {
      'ts-jest': {
        tsconfig: 'tsconfig.json',
      },
    },
    setupFilesAfterEnv: ['<rootDir>/test/redis/redis.mock.ts'],
  };
  