jest.mock('ioredis', () => {
    return {
        __esModule: true,
        default: jest.fn().mockImplementation(() => ({
        get: jest.fn(),
        set: jest.fn(),
        del: jest.fn(),
        })),
    };
});